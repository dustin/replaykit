// Package replay provides a framework for replaying a timestamped
// sequence of events at a relative time.
package replay

import (
	"log"
	"time"
)

// An Event represents individual event found in the source.
type Event interface {
	TS() time.Time // The time this event occurred.
}

// A Source of events (e.g. log reader, etc...)
type Source interface {
	// Next event from the log, nil if there are no more
	Next() Event
}

// Action to process on each event.
type Action interface {
	// Process the event.
	Process(ev Event)
}

// Replay is the primary replay type.  Build a new one with New.
type Replay struct {
	timeScale  float64
	firstEvent time.Time
	realStart  time.Time

	now   func() time.Time
	sleep func(time.Duration)
}

func (r *Replay) timeOffset(eventTime time.Time) time.Duration {
	now := r.now()
	eventElapsed := eventTime.Sub(r.firstEvent)
	localElapsed := time.Duration(float64(now.Sub(r.realStart)) * r.timeScale)

	return time.Duration(float64(eventElapsed-localElapsed) / r.timeScale)

}

func (r *Replay) syncTime(eventTime time.Time) {
	toSleep := r.timeOffset(eventTime)
	if toSleep > 0 {
		r.sleep(toSleep)
	}
}

type functionAction func(Event)

// Process the event.
func (f functionAction) Process(ev Event) {
	f(ev)
}

// FunctionAction wraps a function as an Action.
func FunctionAction(f func(Event)) Action {
	return functionAction(f)
}

type functionSource func() Event

func (f functionSource) Next() Event {
	return f()
}

// FunctionSource creates a source from a function.
func FunctionSource(f func() Event) Source {
	return functionSource(f)
}

// New creates a new Replay with time scaled to the given amount.
//
// scale should be > 0
func New(scale float64) *Replay {
	if scale <= 0 {
		log.Panic("Timescale must be > 0")
	}
	return &Replay{timeScale: scale, now: time.Now, sleep: time.Sleep}
}

// Run the replay.
//
// Returns the amount of time we were "off" of the target.
func (r *Replay) Run(s Source, action Action) time.Duration {
	event := s.Next()
	if event == nil {
		return time.Duration(0)
	}

	r.realStart = r.now()
	r.firstEvent = event.TS()
	eventTime := r.firstEvent

	for ; event != nil; event = s.Next() {

		action.Process(event)

		eventTime = event.TS()
		r.syncTime(eventTime)
	}

	return r.timeOffset(eventTime)
}
