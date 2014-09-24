package replay

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	New(1)
	panicked := false
	func() {
		defer func() { panicked = recover() != nil }()
		r := New(0)
		t.Errorf("Failed to panic at 0, got %#v", r)
	}()
	if !panicked {
		t.Errorf("Failed to panick at zero")
	}
	func() {
		defer func() { panicked = recover() != nil }()
		r := New(-1)
		t.Errorf("Failed to panic at -1, got %#v", r)
	}()
	if !panicked {
		t.Errorf("Failed to panick at -1")
	}
}

type fakeTime struct {
	base   time.Time
	passed time.Duration
}

func (f *fakeTime) now() time.Time {
	return f.base.Add(f.passed)
}

func (f *fakeTime) sleep(d time.Duration) {
	f.passed += d
}

type noopAction struct{}

func (noopAction) Process(Event) {}

type dumbEvent time.Time

func (d dumbEvent) TS() time.Time { return time.Time(d) }

type directSource struct {
	events []Event
}

func (d *directSource) Next() Event {
	if len(d.events) == 0 {
		return nil
	}
	rv := d.events[0]
	d.events = d.events[1:]
	return rv
}

func TestRun(t *testing.T) {
	r := New(1)
	tm := &fakeTime{}
	r.now = tm.now
	r.sleep = tm.sleep

	base := time.Now()
	s := &directSource{
		events: []Event{
			dumbEvent(base.Add(5 * time.Second)),
			dumbEvent(base.Add(6 * time.Second)),
			dumbEvent(base.Add(9 * time.Second)),
			dumbEvent(base.Add(13 * time.Second)),
		}}

	off := r.Run(s, noopAction{})
	if off != 0 {
		t.Errorf("Expected to be off by 0, was off by %v", off)
	}
	if tm.passed != (8 * time.Second) {
		t.Errorf("Expected to take 8 seconds, took %v", tm.passed)
	}
}
