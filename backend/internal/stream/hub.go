package stream

import (
	"sync"
	"sync/atomic"
	"time"
)

// HubEvent is delivered to subscribers via channels.
type HubEvent struct {
	Seq  uint64
	Data string // raw SSE line, e.g. `data: {"type":"text",...}`
	Done bool   // true when the stream is closing
}

// subscriber wraps a per-subscriber buffered channel.
type subscriber struct {
	ch chan HubEvent
}

// RuntimeStream tracks a single active stream with its own seq counter.
type RuntimeStream struct {
	seq         atomic.Uint64
	mu          sync.Mutex
	subscribers map[*subscriber]struct{}
	closed      bool
}

// RuntimeHub is a singleton in-memory pub/sub hub.
// Key format: "sessionID:messageID" (same as Redis stream key suffix).
type RuntimeHub struct {
	mu         sync.RWMutex
	streams    map[string]*RuntimeStream
	closedKeys map[string]struct{} // keys that have been closed; prevents re-creation
}

// Hub is the global RuntimeHub instance.
var Hub = &RuntimeHub{
	streams:    make(map[string]*RuntimeStream),
	closedKeys: make(map[string]struct{}),
}

func (h *RuntimeHub) getOrCreateStream(key string) *RuntimeStream {
	h.mu.RLock()
	s, ok := h.streams[key]
	_, closed := h.closedKeys[key]
	h.mu.RUnlock()
	if ok {
		return s
	}
	if closed {
		return nil // stream was closed; do not re-create
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	// double-check after acquiring write lock
	s, ok = h.streams[key]
	if ok {
		return s
	}
	if _, closed = h.closedKeys[key]; closed {
		return nil
	}
	s = &RuntimeStream{
		subscribers: make(map[*subscriber]struct{}),
	}
	h.streams[key] = s
	return s
}

// Publish sends an event to all subscribers of the given stream key.
// If the stream does not exist, it is created. Publish is non-blocking;
// if a subscriber's buffer is full, the event is silently dropped.
// If the stream key was previously closed, the event is silently discarded.
func (h *RuntimeHub) Publish(key, data string) {
	s := h.getOrCreateStream(key)
	if s == nil {
		return // stream was already closed; drop event
	}
	seq := s.seq.Add(1)
	evt := HubEvent{Seq: seq, Data: data}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	for sub := range s.subscribers {
		select {
		case sub.ch <- evt:
		default:
			// buffer full — drop oldest: drain one, push new
			select {
			case <-sub.ch:
			default:
			}
			select {
			case sub.ch <- evt:
			default:
			}
		}
	}
}

// Subscribe returns a channel for consuming events and the current sequence number.
// The caller uses currentSeq to replay any Redis gap before consuming live events.
// Returns nil channel if the stream key has been closed.
func (h *RuntimeHub) Subscribe(key string) (<-chan HubEvent, uint64) {
	s := h.getOrCreateStream(key)
	if s == nil {
		return nil, 0
	}

	sub := &subscriber{
		ch: make(chan HubEvent, 1024),
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	currentSeq := s.seq.Load()
	if !s.closed {
		s.subscribers[sub] = struct{}{}
	}
	return sub.ch, currentSeq
}

// Unsubscribe removes a subscriber's channel from the stream.
// Call this when the SSE client disconnects to prevent goroutine/channel leaks.
func (h *RuntimeHub) Unsubscribe(key string, ch <-chan HubEvent) {
	h.mu.RLock()
	s, ok := h.streams[key]
	h.mu.RUnlock()
	if !ok {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for sub := range s.subscribers {
		if sub.ch == ch {
			delete(s.subscribers, sub)
			return
		}
	}
}

// Close marks the stream as done, sends a Done event to all subscribers,
// closes their channels, and removes the stream from the hub.
// The key is recorded in closedKeys to prevent re-creation.
func (h *RuntimeHub) Close(key string) {
	h.mu.Lock()
	s, ok := h.streams[key]
	if ok {
		delete(h.streams, key)
	}
	h.closedKeys[key] = struct{}{} // mark as finalized
	h.mu.Unlock()

	if !ok {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	for sub := range s.subscribers {
		// send Done sentinel
		select {
		case sub.ch <- HubEvent{Done: true}:
		default:
		}
		close(sub.ch)
		delete(s.subscribers, sub)
	}
}

// StartClosedKeysCleanup launches a background goroutine that periodically
// resets the closedKeys map. Entries only need to exist long enough to prevent
// re-creation during active streaming; after 10 minutes they are irrelevant.
func (h *RuntimeHub) StartClosedKeysCleanup() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			h.mu.Lock()
			h.closedKeys = make(map[string]struct{})
			h.mu.Unlock()
		}
	}()
}
