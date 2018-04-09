package time_intervals

import (
	"testing"
	"time"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func Test_SubstractBlockedIntervals(t *testing.T) {

	const (
		as  = 1 // as = start of interval A in the docs image
		ss  = 2
		bs  = 3
		ae  = 4 // ae = end of interval A in the docs image, so we can visualise the results
		se  = 5
		ls  = 6
		be  = 7
		le  = 8
		cs  = 9
		ms  = 10
		me  = 11
		ns  = 12
		ce  = 13
		ds  = 14
		os  = 15
		de  = 16
		ne  = 17
		oe  = 18
		es  = 19
		o2s = 20
		fs  = 21
		qs  = 21 // Q starts at same time as F
		qe  = 23
		o2e = 24
		fe  = 25
		ee  = 26
		ps  = 27
		pe  = 28
	)

	A := testInterval(as, ae)
	B := testInterval(bs, be)
	C := testInterval(cs, ce)
	D := testInterval(ds, de)
	E := testInterval(es, ee)
	F := testInterval(fs, fe)

	S := testInterval(ss, se)
	L := testInterval(ls, le)
	M := testInterval(ms, me)
	N := testInterval(ns, ne)
	O := testInterval(os, oe)
	O2 := testInterval(o2s, o2e)
	Q := testInterval(qs, qe)
	P := testInterval(ps, pe)

	available := []Interval{A, B, C, D, E, F}
	blocked := []Interval{S, L, M, N, O, O2, Q, P}

	results := SubstractBlockedIntervals(available, blocked)
	fmt.Printf("Received back %d intervals\n", len(results))
	for _, r := range results {
		printIntervalInMinutes(r)
	}

	assert.Equal(t, 6, len(results), "Expected exactly 6 ordered intervals back")

	// expected intervals:
	// 1 -> 2
	// 5 -> 6
	// 9 -> 10
	// 11 -> 12
	// 19 -> 20
	// 24 -> 26

	assert.Equal(t, 1, results[0].Start.Minute())
	assert.Equal(t, 2, results[0].End.Minute())

	assert.Equal(t, 5, results[1].Start.Minute())
	assert.Equal(t, 6, results[1].End.Minute())

	assert.Equal(t, 9, results[2].Start.Minute())
	assert.Equal(t, 10, results[2].End.Minute())

	assert.Equal(t, 11, results[3].Start.Minute())
	assert.Equal(t, 12, results[3].End.Minute())

	assert.Equal(t, 19, results[4].Start.Minute())
	assert.Equal(t, 20, results[4].End.Minute())

	assert.Equal(t, 24, results[5].Start.Minute())
	assert.Equal(t, 26, results[5].End.Minute())

}

func testInterval(minutesStart int, minutesEnd int) Interval {
	baseTime := time.Date(2018, 4, 7, 0, 0, 0, 0, time.UTC)
	return Interval{
		Start: baseTime.Add(time.Duration(minutesStart) * time.Minute),
		End:   baseTime.Add(time.Duration(minutesEnd) * time.Minute),
	}
}

func onlyIntervalMinutes(i Interval) (int, int) {
	return i.Start.Minute(), i.End.Minute()
}

func printIntervalInMinutes(i Interval) {
	s, e := onlyIntervalMinutes(i)
	fmt.Printf("%d -> %d\n", s, e)
}

func Test_SubstractBlockedIntervals_NoAvailableIntervals(t *testing.T) {
	// available: /empty/
	// blocked:     BBB

	B := testInterval(1, 2)
	results := SubstractBlockedIntervals([]Interval{}, []Interval{B})
	assert.Equal(t, 0, len(results))
}

func Test_SubstractBlockedIntervals_CompletelyOverllaping_Single(t *testing.T) {
	// available  AAA
	// blocked    BBB

	A := testInterval(1, 2)
	B := testInterval(1, 2)
	results := SubstractBlockedIntervals([]Interval{A}, []Interval{B})
	assert.Equal(t, 0, len(results), fmt.Sprintf("result: %v", results))
}

func Test_SubstractBlockedIntervals_CompletelyOverllaping_Multiple(t *testing.T) {
	// available:  AAA CC
	// blocked:    BBB CC

	A := testInterval(1, 2)
	B := testInterval(1, 2)

	C := testInterval(5, 7)
	D := testInterval(5, 7)
	results := SubstractBlockedIntervals([]Interval{A, C}, []Interval{D, B}) // order does not matter
	assert.Equal(t, 0, len(results), fmt.Sprintf("result: %v", results))
}

func Test_SubstractBlockedIntervals_MergesAvailableIntervals_1(t *testing.T) {
	// following intervals are all available, whe check if they get merged
	//  AAA     DD
	//    BBB CCCCC

	A := testInterval(1, 3)
	B := testInterval(2, 4)
	C := testInterval(5, 8)
	D := testInterval(6, 7)

	results := SubstractBlockedIntervals([]Interval{D, B, A, C}, []Interval{}) // order does not matter
	assert.Equal(t, 2, len(results), fmt.Sprintf("result: %v", results))

	assert.Equal(t, 1, results[0].Start.Minute())
	assert.Equal(t, 4, results[0].End.Minute())

	assert.Equal(t, 5, results[1].Start.Minute())
	assert.Equal(t, 8, results[1].End.Minute())
}

func Test_SubstractBlockedIntervals_MergesAvailableIntervals_2(t *testing.T) {
	// following intervals are all available, whe check if they get merged
	//  AAA     DD  EEE
	//    BBB CCCCCCCCCCC

	A := testInterval(1, 3)
	B := testInterval(2, 4)
	C := testInterval(5, 15)
	D := testInterval(6, 7)
	E := testInterval(10, 12)

	results := SubstractBlockedIntervals([]Interval{D, E, B, C, A}, []Interval{}) // order does not matter
	assert.Equal(t, 2, len(results), fmt.Sprintf("result: %v", results))

	assert.Equal(t, 1, results[0].Start.Minute())
	assert.Equal(t, 4, results[0].End.Minute())

	assert.Equal(t, 5, results[1].Start.Minute())
	assert.Equal(t, 15, results[1].End.Minute())
}
