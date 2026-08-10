package display

import "sync"

// Hub broadcasts the latest display frame without allowing slow subscribers
// to block the renderer. Each subscriber retains at most one pending frame.
type Hub struct {
	mu          sync.Mutex
	latest      Frame
	hasLatest   bool
	subscribers map[chan Frame]struct{}
}

func NewHub() *Hub {
	return &Hub{subscribers: make(map[chan Frame]struct{})}
}

func (h *Hub) Publish(frame Frame) {
	frame = frame.clone()

	h.mu.Lock()
	defer h.mu.Unlock()

	h.latest = frame
	h.hasLatest = true
	for subscriber := range h.subscribers {
		queued := frame.clone()
		select {
		case subscriber <- queued:
		default:
			select {
			case <-subscriber:
			default:
			}
			subscriber <- queued
		}
	}
}

// Subscribe returns a stream that immediately contains the latest frame, when
// one exists, and a function that removes the subscription.
func (h *Hub) Subscribe() (<-chan Frame, func()) {
	subscriber := make(chan Frame, 1)

	h.mu.Lock()
	h.subscribers[subscriber] = struct{}{}
	if h.hasLatest {
		subscriber <- h.latest.clone()
	}
	h.mu.Unlock()

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			h.mu.Lock()
			delete(h.subscribers, subscriber)
			close(subscriber)
			h.mu.Unlock()
		})
	}
	return subscriber, unsubscribe
}
