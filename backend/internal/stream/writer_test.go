package stream

import (
	"testing"
	"time"
)

// Note: PublishErrorAndFail requires db.GetDB() and redis to be initialized.
// We test Hub behavior separately since the DB dependency is not mockable in unit tests.

func TestHub_ClosePreventsRecreation(t *testing.T) {
	Hub = &RuntimeHub{
		streams:    make(map[string]*RuntimeStream),
		closedKeys: make(map[string]struct{}),
	}

	key := "session:msg-close-test"

	// Publish to create stream
	Hub.Publish(key, "data: test")
	Hub.Close(key)

	// Publish again — should be silently dropped (no re-creation)
	Hub.Publish(key, "data: after-close")

	Hub.mu.RLock()
	_, exists := Hub.streams[key]
	_, closed := Hub.closedKeys[key]
	Hub.mu.RUnlock()

	if exists {
		t.Error("stream should not exist after Close")
	}
	if !closed {
		t.Error("key should be in closedKeys after Close")
	}
}

func TestHub_SubscribeReturnsNilAfterClose(t *testing.T) {
	Hub = &RuntimeHub{
		streams:    make(map[string]*RuntimeStream),
		closedKeys: make(map[string]struct{}),
	}

	key := "session:msg-sub-nil"

	// Close without ever publishing — mark key as closed
	Hub.Close(key)

	ch, seq := Hub.Subscribe(key)
	if ch != nil {
		t.Error("expected nil channel when subscribing to a closed key")
	}
	if seq != 0 {
		t.Errorf("expected seq 0 for closed key, got %d", seq)
	}
}

func TestHub_UnsubscribeRemovesSubscriber(t *testing.T) {
	Hub = &RuntimeHub{
		streams:    make(map[string]*RuntimeStream),
		closedKeys: make(map[string]struct{}),
	}

	key := "session:msg-unsub"

	// Create stream and subscribe
	Hub.Publish(key, "data: init")
	ch, _ := Hub.Subscribe(key)

	// Verify subscriber exists
	Hub.mu.RLock()
	s := Hub.streams[key]
	Hub.mu.RUnlock()
	s.mu.Lock()
	count := len(s.subscribers)
	s.mu.Unlock()
	if count != 1 {
		t.Fatalf("expected 1 subscriber, got %d", count)
	}

	// Unsubscribe
	Hub.Unsubscribe(key, ch)

	// Verify subscriber removed
	s.mu.Lock()
	count = len(s.subscribers)
	s.mu.Unlock()
	if count != 0 {
		t.Errorf("expected 0 subscribers after Unsubscribe, got %d", count)
	}
}

func TestHub_UnsubscribeOnNonexistentStream(t *testing.T) {
	Hub = &RuntimeHub{
		streams:    make(map[string]*RuntimeStream),
		closedKeys: make(map[string]struct{}),
	}

	// Should not panic
	ch := make(chan HubEvent, 10)
	Hub.Unsubscribe("nonexistent:key", ch)
}

func TestHub_PublishDropOnClosedStream(t *testing.T) {
	Hub = &RuntimeHub{
		streams:    make(map[string]*RuntimeStream),
		closedKeys: make(map[string]struct{}),
	}

	key := "session:msg-drop"

	// Create and close
	Hub.Publish(key, "data: before")
	Hub.Close(key)

	// This should not create a new stream
	Hub.Publish(key, "data: after-close")

	Hub.mu.RLock()
	_, exists := Hub.streams[key]
	Hub.mu.RUnlock()
	if exists {
		t.Error("stream should not be recreated after Close")
	}
}

func TestHub_StartClosedKeysCleanup(t *testing.T) {
	Hub = &RuntimeHub{
		streams:    make(map[string]*RuntimeStream),
		closedKeys: make(map[string]struct{}),
	}

	// Add some closed keys
	Hub.Close("key1")
	Hub.Close("key2")

	Hub.mu.RLock()
	count := len(Hub.closedKeys)
	Hub.mu.RUnlock()
	if count != 2 {
		t.Fatalf("expected 2 closedKeys, got %d", count)
	}

	// Start cleanup (runs every 10min in prod, but we just test it doesn't panic)
	// We won't wait 10 minutes; just verify it starts without error
	done := make(chan struct{})
	go func() {
		Hub.StartClosedKeysCleanup()
		close(done)
	}()

	// Give it a moment to start
	time.Sleep(50 * time.Millisecond)

	Hub.mu.RLock()
	count = len(Hub.closedKeys)
	Hub.mu.RUnlock()
	// Keys should still exist (cleanup hasn't run yet — 10 min interval)
	if count != 2 {
		t.Logf("closedKeys count = %d (cleanup may have run)", count)
	}
}
