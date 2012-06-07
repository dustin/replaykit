// Timed event replay toolkit.
package replay

import (
	"log"
	"time"
)

// An individual event found in the source.
type Event interface {
	TS() time.Time // The time this event occurred.
}

// A source of events (e.g. log reader, etc...)
type Source interface {
	// Next event from the log, nil if there are no more
	Next() Event
}

// Action to process on each event.
type Action interface {
	// Process the event.
	Process(ev Event)
}

// The replayer.  Build a new one with New.
type Replay struct {
	timeScale  float64
	firstEvent time.Time
	realStart  time.Time
}

func (r *Replay) timeOffset(eventTime time.Time) time.Duration {
	now := time.Now()
	eventElapsed := eventTime.Sub(r.firstEvent)
	localElapsed := time.Duration(float64(now.Sub(r.realStart)) * r.timeScale)

	return time.Duration(float64(eventElapsed-localElapsed) / r.timeScale)

}

func (r *Replay) syncTime(eventTime time.Time) {
	toSleep := r.timeOffset(eventTime)
	if toSleep > 0 {
		time.Sleep(toSleep)
	}
}

// Build a new Replayer with time scaled to the given amount.
//
// scale should be > 0
func New(scale float64) *Replay {
	if scale <= 0 {
		log.Panic("Timescale must be > 0")
	}
	return &Replay{timeScale: scale}
}

// Run the replay.
//
// Returns the amount of time we were "off" of the target.
func (r *Replay) Run(s Source, action Action) time.Duration {

	if r.timeScale <= 0 {
		log.Panic("Timescale must be > 0")
	}

	event := s.Next()
	if event == nil {
		return time.Duration(0)
	}

	r.realStart = time.Now()
	r.firstEvent = event.TS()
	eventTime := r.firstEvent

	for ; event != nil; event = s.Next() {

		action.Process(event)

		eventTime = event.TS()
		r.syncTime(eventTime)
	}

	return r.timeOffset(eventTime)
}
