package replay

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"time"
)

// exampleData for the example source.
//
// Imagine you have an arbitrarily long stream in this form.
const exampleData = `2014-09-24T19:47:32-07:00 first
2014-09-24T19:47:42-07:00 second
2014-09-24T19:48:13-07:00 third`

// exampleEvent represents a single event extracted from that data.
type exampleEvent struct {
	theTime   time.Time
	eventName string
}

// TS just returns the parsed time from the event.
func (t exampleEvent) TS() time.Time { return t.theTime }

// exampleSrc contains a bufio.Scanner with which it parses the stream.
type exampleSrc struct {
	r *bufio.Scanner
}

// Next parses the next event out of the stream and emits it.
//
// This is a super simple example made for this arbitrary test, but
// it's basically how most of them look.  You build a parser for
// whatever data you've got that contains timestamped events, and work
// with one at a time.
func (e exampleSrc) Next() Event {
	if !e.r.Scan() {
		return nil
	}
	parts := strings.Split(e.r.Text(), " ")
	t, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		log.Printf("Failed to parse %q: %v", parts[0], err)
		return nil
	}

	return exampleEvent{t, parts[1]}
}

// Example is a complete example of processing data.
//
// A couple things to note about this:
//  1. It uses fake time (so it runs instantly).
//  2. As such, it introduces fake overhead.
//  3. It's running at 10x realtime, you could also run at 0.1x realtime, but
//     it's harder to introduce noticable simulated overhead for the processing
//     work.
func Example() {
	src := exampleSrc{bufio.NewScanner(strings.NewReader(exampleData))}
	r, tm := useFakeTime(New(10))
	off := r.Run(src, FunctionAction(func(e Event) {
		myEvent := e.(exampleEvent)
		// Predictable overhead time
		tm.sleep(time.Second * time.Duration(len(myEvent.eventName)))
		fmt.Printf("Processing event named %q at %v\n", myEvent.eventName, myEvent.theTime)
	}))
	fmt.Printf("Took %v with final entry off by %v\n", tm.passed, off)
	// Output:
	// Processing event named "first" at 2014-09-24 19:47:32 -0700 PDT
	// Processing event named "second" at 2014-09-24 19:47:42 -0700 PDT
	// Processing event named "third" at 2014-09-24 19:48:13 -0700 PDT
	// Took 16s with final entry off by -11.9s
}
