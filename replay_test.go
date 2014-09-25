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

func useFakeTime(r *Replay) (*Replay, *fakeTime) {
	tm := &fakeTime{}
	r.now = tm.now
	r.sleep = tm.sleep
	return r, tm
}

type dumbEvent time.Time

func (d dumbEvent) TS() time.Time { return time.Time(d) }

func genEvents() []Event {
	base := time.Now()
	return []Event{
		dumbEvent(base.Add(5 * time.Second)),
		dumbEvent(base.Add(6 * time.Second)),
		dumbEvent(base.Add(9 * time.Second)),
		dumbEvent(base.Add(13 * time.Second)),
	}
}

var noopAction = FunctionAction(func(Event) {})

func TestRun(t *testing.T) {
	r, tm := useFakeTime(New(1))

	off := r.Run(CollectionSource(genEvents()), noopAction)
	if off != 0 {
		t.Errorf("Expected to be off by 0, was off by %v", off)
	}
	if tm.passed != (8 * time.Second) {
		t.Errorf("Expected to take 8 seconds, took %v", tm.passed)
	}
}

func TestRun2x(t *testing.T) {
	r, tm := useFakeTime(New(2))

	off := r.Run(CollectionSource(genEvents()), noopAction)
	if off != 0 {
		t.Errorf("Expected to be off by 0, was off by %v", off)
	}
	if tm.passed != (4 * time.Second) {
		t.Errorf("Expected to take 8 seconds, took %v", tm.passed)
	}
}

func TestRunHalfx(t *testing.T) {
	r, tm := useFakeTime(New(0.5))

	off := r.Run(CollectionSource(genEvents()), noopAction)
	if off != 0 {
		t.Errorf("Expected to be off by 0, was off by %v", off)
	}
	if tm.passed != (16 * time.Second) {
		t.Errorf("Expected to take 8 seconds, took %v", tm.passed)
	}
}

func TestRunNil(t *testing.T) {
	r, _ := useFakeTime(New(1))
	off := r.Run(CollectionSource(nil), FunctionAction(func(Event) {}))
	if off != 0 {
		t.Errorf("Expected nil input to run with 0 off, got %v", off)
	}
}
