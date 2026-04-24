package env

import (
	"errors"
	"sync"
	"time"
)

// ChangeKind describes the type of change observed on a Set.
type ChangeKind string

const (
	ChangeKindPut    ChangeKind = "put"
	ChangeKindDelete ChangeKind = "delete"
)

// ChangeEvent is emitted by a Watcher whenever a key in the watched Set
// is created, updated, or deleted.
type ChangeEvent struct {
	Key      string
	Kind     ChangeKind
	OldValue string // empty when Kind == ChangeKindPut and key is new
	NewValue string // empty when Kind == ChangeKindDelete
	At       time.Time
}

// Watcher polls a Set at a configurable interval and emits ChangeEvents
// to registered subscribers whenever the contents change.
type Watcher struct {
	mu       sync.Mutex
	set      *Set
	interval time.Duration
	last     map[string]string
	subs     []chan ChangeEvent
	stop     chan struct{}
	stopped  chan struct{}
}

// NewWatcher creates a Watcher that observes set at the given poll interval.
// The interval must be positive and set must not be nil.
func NewWatcher(set *Set, interval time.Duration) (*Watcher, error) {
	if set == nil {
		return nil, errors.New("watch: set must not be nil")
	}
	if interval <= 0 {
		return nil, errors.New("watch: interval must be positive")
	}
	snapshot, err := TakeSnapshot(set)
	if err != nil {
		return nil, err
	}
	return &Watcher{
		set:      set,
		interval: interval,
		last:     snapshot.Vars,
		stop:     make(chan struct{}),
		stopped:  make(chan struct{}),
	}, nil
}

// Subscribe returns a channel that receives ChangeEvents while the Watcher
// is running. The channel is buffered with capacity 64 to avoid blocking
// the poll loop on slow consumers.
func (w *Watcher) Subscribe() <-chan ChangeEvent {
	ch := make(chan ChangeEvent, 64)
	w.mu.Lock()
	w.subs = append(w.subs, ch)
	w.mu.Unlock()
	return ch
}

// Start begins the polling loop in a background goroutine.
// Calling Start on an already-running Watcher is a no-op.
func (w *Watcher) Start() {
	go w.loop()
}

// Stop signals the polling loop to exit and waits for it to finish.
func (w *Watcher) Stop() {
	close(w.stop)
	<-w.stopped
	w.mu.Lock()
	for _, ch := range w.subs {
		close(ch)
	}
	w.subs = nil
	w.mu.Unlock()
}

func (w *Watcher) loop() {
	defer close(w.stopped)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-w.stop:
			return
		case <-ticker.C:
			w.poll()
		}
	}
}

func (w *Watcher) poll() {
	snap, err := TakeSnapshot(w.set)
	if err != nil {
		return
	}
	now := time.Now()
	current := snap.Vars
	w.mu.Lock()
	defer w.mu.Unlock()

	// Detect puts (new or updated keys).
	for k, newVal := range current {
		oldVal, existed := w.last[k]
		if !existed || oldVal != newVal {
			ev := ChangeEvent{Key: k, Kind: ChangeKindPut, OldValue: oldVal, NewValue: newVal, At: now}
			w.broadcast(ev)
		}
	}
	// Detect deletes.
	for k, oldVal := range w.last {
		if _, exists := current[k]; !exists {
			ev := ChangeEvent{Key: k, Kind: ChangeKindDelete, OldValue: oldVal, At: now}
			w.broadcast(ev)
		}
	}
	w.last = current
}

// broadcast sends ev to all subscribers; drops the event if a channel is full.
func (w *Watcher) broadcast(ev ChangeEvent) {
	for _, ch := range w.subs {
		select {
		case ch <- ev:
		default:
		}
	}
}
