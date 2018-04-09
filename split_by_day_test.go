package time_intervals

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"time"
)

// uses april 10 as base date, so day 7 => 17 and it is easier to follow
var baseTime = time.Date(2018, 4, 10, 0, 0, 0, 0, time.UTC)

func Test_IntervalsByDay(t *testing.T) {
	// day0     day1	day2	day3	day4	day5	day6	day7
	// |		|		|		|		|		|		|		|
	//   AAA   BBB  C     DDDDDDDDDDDDDDDDDDD EEEEEEEEEE  F G H JJJ

	A := testDHInterval(0, 3, 0, 7)
	B := testDHInterval(0, 17, 1, 5)
	C := testDHInterval(1, 11, 1, 13)
	D := testDHInterval(2, 8, 4, 9)
	E := testDHInterval(4, 21, 5, 24)
	E.End = E.End.Add(-time.Millisecond)

	F := testDHInterval(6, 8, 6, 9)
	G := testDHInterval(6, 10, 6, 11)
	H := testDHInterval(6, 13, 6, 15)
	J := testDHInterval(7, 0, 7, 9)

	results := IntervalsByDay([]Interval{A, B, C, D, E, F, G, H, J}) // order DOES matter
	assert.Equal(t, 8, len(results), fmt.Sprintf("Expected exactly 8 keys coresponding to days 0->8, but days were: %v", getKeys(results)))

	// Day 0
	assert.Equal(t, 2, len(results[NormalizeDate(A.Start)]), "Day 0 should have 2 intervals, but results were: %+v", results[NormalizeDate(A.Start)])

	// A interval
	e := testDHInterval(0, 3, 0, 7)
	r := results[NormalizeDate(A.Start)][0]
	assert.True(t, intervalsDiff(e, r) == "", intervalsDiff(e, r))

	// B start
	e = testDHInterval(0, 17, 1, 0)
	e.End = e.End.Add(-time.Millisecond)
	r = results[NormalizeDate(A.Start)][1]
	assert.True(t, intervalsDiff(e, r) == "", intervalsDiff(e, r))

	// Day 1
	assert.Equal(t, 2, len(results[NormalizeDate(A.Start)]), "Day 1 should have 2 intervals, but results were: %+v", results[NormalizeDate(B.Start)])

	// B continuation
	e = testDHInterval(1, 0, 1, 5)
	r = results[NormalizeDate(B.End)][0]
	assert.True(t, intervalsDiff(e, r) == "", intervalsDiff(e, r))

	// C
	e = testDHInterval(1, 11, 1, 13)
	r = results[NormalizeDate(B.End)][1]
	assert.True(t, intervalsDiff(e, r) == "", intervalsDiff(e, r))

	// Day 2
	assert.Equal(t, 1, len(results[NormalizeDate(D.Start)]), "Day 2 should have 1 interval, but results were: %+v", results[NormalizeDate(D.Start)])

	// D start
	e = testDHInterval(2, 8, 3, 0)
	e.End = e.End.Add(-time.Millisecond)
	r = results[NormalizeDate(D.Start)][0]
	assert.True(t, intervalsDiff(e, r) == "", intervalsDiff(e, r))

	// Day 3
	assert.Equal(t, 1, len(results[NormalizeDate(D.Start).AddDate(0, 0, 1)]), "Day 3 should have 1 interval, but results were: %+v", results[NormalizeDate(D.Start).AddDate(0, 0, 1)])
	e = testDHInterval(3, 0, 4, 0)
	e.End = e.End.Add(-time.Millisecond)
	r = results[NormalizeDate(D.Start).AddDate(0, 0, 1)][0]
	assert.True(t, intervalsDiff(e, r) == "", intervalsDiff(e, r))

	// Day 4
	assert.Equal(t, 2, len(results[NormalizeDate(E.Start)]), "Day 4 should have 2 intervals, but results were: %+v", results[NormalizeDate(E.Start)])
	// D ending
	e = testDHInterval(4, 0, 4, 9)
	r = results[NormalizeDate(D.End)][0]
	assert.True(t, intervalsDiff(e, r) == "", intervalsDiff(e, r))

	// E start
	e = testDHInterval(4, 21, 5, 0)
	e.End = e.End.Add(-time.Millisecond)
	r = results[NormalizeDate(D.End)][1]
	assert.True(t, intervalsDiff(e, r) == "", intervalsDiff(e, r))

	// Day 5
	assert.Equal(t, 1, len(results[NormalizeDate(E.Start).AddDate(0, 0, 1)]), "Day 5 should have 1 interval, but results were: %+v", results[NormalizeDate(E.Start).AddDate(0, 0, 1)])
	// E ending
	e = testDHInterval(5, 0, 6, 0)
	e.End = e.End.Add(-time.Millisecond)
	r = results[NormalizeDate(E.Start.AddDate(0, 0, 1))][0]
	assert.True(t, intervalsDiff(e, r) == "", intervalsDiff(e, r))

	// Day 6
	assert.Equal(t, 3, len(results[NormalizeDate(F.Start)]), "Day 6 should have 3 intervals, but results were: %+v", results[NormalizeDate(F.Start)])

	// Day 7
	assert.Equal(t, 1, len(results[NormalizeDate(J.Start)]), "Day 7 should have 3 intervals, but results were: %+v", results[NormalizeDate(J.Start)])
	// J
	e = testDHInterval(7, 0, 7, 9)
	r = results[NormalizeDate(J.Start)][0]
	assert.True(t, intervalsDiff(e, r) == "", intervalsDiff(e, r))
}

func testDHInterval(dStart, hStart, dEnd, hEnd int) Interval {

	return Interval{
		Start: baseTime.AddDate(0, 0, dStart).Add(time.Duration(hStart) * time.Hour),
		End:   baseTime.AddDate(0, 0, dEnd).Add(time.Duration(hEnd) * time.Hour),
	}
}

func intervalsDiff(a, b Interval) string {
	if a.Start != b.Start {
		return fmt.Sprintf("start time not equal: %v <> %v", a.Start, b.Start)
	}
	if a.End != b.End {
		return fmt.Sprintf("end time not equal: %v <> %v", a.End, b.End)
	}
	return ""
}

func getKeys(m map[time.Time][]Interval) []time.Time {
	res := []time.Time{}
	for k, _ := range m {
		res = append(res, k)
	}
	return res
}

func Test_IntervalsForEachDayInRange(t *testing.T) {
	// day0     day1	day2	day3	day4	day5	day6	day7
	// |		|		|		|		|		|		|		|
	//   AAA   BBB  C     DDDDDDDDDDDDDDDDDDD             F G H JJJ

	A := testDHInterval(0, 3, 0, 7)
	B := testDHInterval(0, 17, 1, 5)
	C := testDHInterval(1, 11, 1, 13)
	D := testDHInterval(2, 8, 4, 9)

	// E is missing

	F := testDHInterval(6, 8, 6, 9)
	G := testDHInterval(6, 10, 6, 11)
	H := testDHInterval(6, 13, 6, 15)
	J := testDHInterval(7, 0, 7, 9)

	d, err := IntervalsForEachDayInRange([]Interval{A, B, C, D, F, G, H, J}, A.Start.AddDate(0, 0, -3), J.End.AddDate(0, 0, 1))
	assert.NoError(t, err)

	assert.Equal(t, 12, len(d), "Expected all days in the range Apr 7 -> Apr 17 inclusive")

	assert.Equal(t, NormalizeDate(A.Start).AddDate(0, 0, -3), d[0].Date)
	assert.Equal(t, 0, d[0].CountSinceFirst)
	assert.Equal(t, 0, len(d[0].OrderedDisjunctIntervals))

	assert.Equal(t, NormalizeDate(A.Start).AddDate(0, 0, -2), d[1].Date)
	assert.Equal(t, 1, d[1].CountSinceFirst)
	assert.Equal(t, 0, len(d[1].OrderedDisjunctIntervals))

	assert.Equal(t, NormalizeDate(A.Start).AddDate(0, 0, -1), d[2].Date)
	assert.Equal(t, 2, d[2].CountSinceFirst)
	assert.Equal(t, 0, len(d[2].OrderedDisjunctIntervals))

	assert.Equal(t, NormalizeDate(A.Start).AddDate(0, 0, 0), d[3].Date)
	assert.Equal(t, 3, d[3].CountSinceFirst)
	assert.Equal(t, 2, len(d[3].OrderedDisjunctIntervals))

	assert.Equal(t, NormalizeDate(A.Start).AddDate(0, 0, 1), d[4].Date)
	assert.Equal(t, 4, d[4].CountSinceFirst)
	assert.Equal(t, 2, len(d[4].OrderedDisjunctIntervals))

	assert.Equal(t, NormalizeDate(A.Start).AddDate(0, 0, 2), d[5].Date)
	assert.Equal(t, 5, d[5].CountSinceFirst)
	assert.Equal(t, 1, len(d[5].OrderedDisjunctIntervals))

	assert.Equal(t, NormalizeDate(A.Start).AddDate(0, 0, 4), d[7].Date)
	assert.Equal(t, 7, d[7].CountSinceFirst)
	assert.Equal(t, 1, len(d[7].OrderedDisjunctIntervals))

	assert.Equal(t, NormalizeDate(A.Start).AddDate(0, 0, 5), d[8].Date)
	assert.Equal(t, 8, d[8].CountSinceFirst)
	assert.Equal(t, 0, len(d[8].OrderedDisjunctIntervals))

	assert.Equal(t, NormalizeDate(A.Start).AddDate(0, 0, 6), d[9].Date)
	assert.Equal(t, 9, d[9].CountSinceFirst)
	assert.Equal(t, 3, len(d[9].OrderedDisjunctIntervals))
}

func Test_SplitInFixedIntervals(t *testing.T) {
	// 8am      9am     10am    11am    12am    1pm     2pm     3pm
	// |		|		|		|		|		|		|		|
	//    AAAAAAAAAAAAAAAAAAAAAAAAAAA  BBBBBBBB     CCCCCCCC
	//    000011112222333344445555666  77778888     99990000

	A := testInterval(8*60+15, 11*60+43)
	B := testInterval(12*60, 13*60)
	C := testInterval(13*60+30, 14*60+40)

	r := SplitInFixedIntervals([]Interval{A, B, C}, 30)

	assert.Equal(t, 10, len(r))

	assert.Equal(t, r[0].Start, A.Start)
	assert.Equal(t, r[0].End, A.Start.Add(time.Duration(30)*time.Minute))

	assert.Equal(t, r[1].Start, A.Start.Add(time.Duration(30)*time.Minute))
	assert.Equal(t, r[1].End, A.Start.Add(time.Duration(60)*time.Minute))

	assert.Equal(t, r[2].Start, A.Start.Add(time.Duration(60)*time.Minute))
	assert.Equal(t, r[2].End, A.Start.Add(time.Duration(90)*time.Minute))

	assert.Equal(t, r[3].Start, A.Start.Add(time.Duration(90)*time.Minute))
	assert.Equal(t, r[3].End, A.Start.Add(time.Duration(120)*time.Minute))

	assert.Equal(t, r[4].Start, A.Start.Add(time.Duration(120)*time.Minute))
	assert.Equal(t, r[4].End, A.Start.Add(time.Duration(150)*time.Minute))

	assert.Equal(t, r[5].Start, A.Start.Add(time.Duration(150)*time.Minute))
	assert.Equal(t, r[5].End, A.Start.Add(time.Duration(180)*time.Minute))

	assert.Equal(t, r[6].Start, B.Start.Add(time.Duration(0)*time.Minute))
	assert.Equal(t, r[6].End, B.Start.Add(time.Duration(30)*time.Minute))
}
