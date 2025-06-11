package service

import (
	"io"
	"log"
	"testing"
)

func TestWorkerCapacity(t *testing.T) {
	logger := log.New(io.Discard, "", 0)
	wIface, err := NewWorker(logger, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	w := wIface.(*worker)
	// reduce size for testing
	w.messages = make(chan Message, 2)
	w.results = make(chan Result, 2)
	w.size = 2
	w.state.Store(STATE_RUNNING)

	current, size := w.Capacity()
	if current != size {
		t.Fatalf("expected capacity %d, got %d", size, current)
	}

	if err := w.AddMessage(Message{topic: "test", payload: []byte("1")}); err != nil {
		t.Fatalf("unexpected add error: %v", err)
	}

	current, _ = w.Capacity()
	if current != size-1 {
		t.Fatalf("expected capacity %d, got %d", size-1, current)
	}

	w.Stop()
}
