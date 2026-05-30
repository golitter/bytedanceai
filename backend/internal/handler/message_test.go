package handler

import (
	"testing"

	"agenthub/backend/internal/model"
)

func TestReverseMessages(t *testing.T) {
	messages := []model.Message{
		{ID: 3, Content: "newest"},
		{ID: 2, Content: "middle"},
		{ID: 1, Content: "oldest"},
	}

	reverseMessages(messages)

	got := []uint{messages[0].ID, messages[1].ID, messages[2].ID}
	want := []uint{1, 2, 3}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("reverseMessages() = %v, want %v", got, want)
		}
	}
}
