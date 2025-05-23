package events

import (
	"sync"
)

// EventHandler is a function that handles domain events.
type EventHandler func(event Event)

// EventDispatcher is responsible for dispatching domain events to registered handlers.
type EventDispatcher interface {
	// Register adds a new event handler for a specific event type.
	Register(eventType EventType, handler EventHandler)

	// Unregister removes an event handler for a specific event type.
	Unregister(eventType EventType, handler EventHandler)

	// Dispatch dispatches an event to all registered handlers.
	Dispatch(event Event)
}

// simpleEventDispatcher is a simple in-memory implementation of EventDispatcher.
type simpleEventDispatcher struct {
	handlers map[EventType][]EventHandler
	mu       sync.RWMutex
}

// NewSimpleEventDispatcher creates a new simple event dispatcher.
func NewSimpleEventDispatcher() EventDispatcher {
	return &simpleEventDispatcher{
		handlers: make(map[EventType][]EventHandler),
	}
}

// Register adds a new event handler for a specific event type.
func (d *simpleEventDispatcher) Register(eventType EventType, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.handlers[eventType]; !ok {
		d.handlers[eventType] = []EventHandler{}
	}

	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Unregister removes an event handler for a specific event type.
func (d *simpleEventDispatcher) Unregister(eventType EventType, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if handlers, ok := d.handlers[eventType]; ok {
		for i, h := range handlers {
			// This is a simplistic approach to comparing functions
			// In a more sophisticated implementation, we might use function wrappers with IDs
			if &h == &handler {
				d.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// Dispatch dispatches an event to all registered handlers.
func (d *simpleEventDispatcher) Dispatch(event Event) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	eventType := event.Type()

	// Dispatch to specific event type handlers
	if handlers, ok := d.handlers[eventType]; ok {
		for _, handler := range handlers {
			// In a production implementation, we might want to:
			// - Execute handlers in a separate goroutine
			// - Add error handling
			// - Add event logging
			// - Add retry mechanisms
			handler(event)
		}
	}

	// Also dispatch to "any" event handlers that want all events
	if allHandlers, ok := d.handlers[""]; ok {
		for _, handler := range allHandlers {
			handler(event)
		}
	}
}
