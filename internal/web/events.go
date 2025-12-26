package web

import "sync"

type eventTopic string

const (
	eventShoppingList eventTopic = "shopping-list"
	eventProductsList eventTopic = "products-list"
)

type eventHub struct {
	mu   sync.Mutex
	subs map[eventTopic]map[*eventSub]struct{}
}

type eventSub struct {
	ch       chan struct{}
	clientID string
}

func newEventHub() *eventHub {
	return &eventHub{
		subs: map[eventTopic]map[*eventSub]struct{}{
			eventShoppingList: {},
			eventProductsList: {},
		},
	}
}

func (h *eventHub) Subscribe(topic eventTopic, clientID string) (<-chan struct{}, func()) {
	sub := &eventSub{
		ch:       make(chan struct{}, 1),
		clientID: clientID,
	}
	h.mu.Lock()
	h.subs[topic][sub] = struct{}{}
	h.mu.Unlock()

	return sub.ch, func() {
		h.mu.Lock()
		delete(h.subs[topic], sub)
		h.mu.Unlock()
		close(sub.ch)
	}
}

func (h *eventHub) Publish(topic eventTopic, excludeClientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for sub := range h.subs[topic] {
		if excludeClientID != "" && sub.clientID == excludeClientID {
			continue
		}
		select {
		case sub.ch <- struct{}{}:
		default:
		}
	}
}
