package time_intervals

import (
	"time"
	"container/heap"
	"fmt"
	"github.com/go-errors/errors"
)

type Interval struct {
	Start time.Time
	End   time.Time
}

func (i *Interval) String() string {
	return fmt.Sprintf("(%v -> %v)[%f minutes]", i.Start, i.End, i.End.Sub(i.Start).Minutes())
}

type IntervalType string

const Available IntervalType = "available"
const Blocked IntervalType = "blocked"

type EndpointType string

const Start EndpointType = "start"
const End EndpointType = "end"

type Endpoint struct {
	IntervalType IntervalType
	EndpointType EndpointType
	Time         time.Time
}

type EndpointsHeap []Endpoint

func (h EndpointsHeap) Len() int { return len(h) }
func (h EndpointsHeap) Less(i, j int) bool {
	d := h[i].Time.Sub(h[j].Time)
	if d < 0 {
		return true
	} else if d > 0 {
		return false
	} else {
		// prefer open intervals first which will merge adjacent intervals inside SubstractBlockedIntervals
		if h[i].EndpointType == Start {
			return true
		}
		if h[j].EndpointType == Start {
			return false
		}
		return true
	}
}
func (h EndpointsHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *EndpointsHeap) Push(x interface{}) {
	*h = append(*h, x.(Endpoint))
}

func (h *EndpointsHeap) Pop() interface{} {
	oldh := *h
	x := oldh[len(oldh)-1]
	newh := oldh[0:len(oldh)-1]
	*h = newh
	return x
}

func SubstractBlockedIntervals(available []Interval, blocked []Interval) []Interval {
	h := buildEndpointsHeap(available, blocked)
	nItems := len(*h)

	results := []Interval{}

	availableOpenIntervals := 0
	blockedOpenIntervals := 0
	currentAvailableIntervalStart := time.Time{}
	for i := 0; i < nItems; i++ {
		e := heap.Pop(h).(Endpoint)

		nextAvailableOpenIntervals, nextBlockedOpenIntervals := getNextCounts(e, availableOpenIntervals, blockedOpenIntervals)
		// fmt.Printf("Loop %d: [endpoint %d] (%d, %d) -> (%d, %d)\n", i, e.Time.Minute(), availableOpenIntervals, blockedOpenIntervals, nextAvailableOpenIntervals, nextBlockedOpenIntervals)
		if blockedOpenIntervals > 0 && nextBlockedOpenIntervals > 0 {
			// do nothing, we were blocked before, we are still blocked
		}
		if blockedOpenIntervals == 0 && nextBlockedOpenIntervals > 0 {
			// we are getting blocked
			if availableOpenIntervals > 0 {
				if currentAvailableIntervalStart.Before(e.Time) {
					results = append(results, Interval{Start: currentAvailableIntervalStart, End: e.Time})
				}
			}
		}
		if blockedOpenIntervals > 0 && nextBlockedOpenIntervals == 0 {
			// we got unblocked
			if availableOpenIntervals > 0 {
				// next available interval should start here
				currentAvailableIntervalStart = e.Time
			}
		}
		if blockedOpenIntervals == 0 && nextBlockedOpenIntervals == 0 {
			// there is no blocking interval
			if availableOpenIntervals > 0 && nextAvailableOpenIntervals == 0 {
				// the current available interval is ending
				if currentAvailableIntervalStart.Before(e.Time) {
					results = append(results, Interval{Start: currentAvailableIntervalStart, End: e.Time})
				}
			}
			if availableOpenIntervals == 0 && nextAvailableOpenIntervals > 0 {
				// a new available interval is starting
				currentAvailableIntervalStart = e.Time
			}
		}

		// assign next values
		availableOpenIntervals, blockedOpenIntervals = nextAvailableOpenIntervals, nextBlockedOpenIntervals
	}
	return results
}

func buildEndpointsHeap(available []Interval, blocked []Interval) *EndpointsHeap {
	endpoints := EndpointsHeap{}
	for _, a := range available {
		endpoints.Push(Endpoint{Available, Start, a.Start})
		endpoints.Push(Endpoint{Available, End, a.End})
	}
	for _, b := range blocked {
		endpoints.Push(Endpoint{Blocked, Start, b.Start})
		endpoints.Push(Endpoint{Blocked, End, b.End})
	}
	heap.Init(&endpoints)
	return &endpoints
}

func getNextCounts(e Endpoint, availableOpenIntervals int, blockedOpenIntervals int) (nextAvailable int, nextBlocked int) {
	nextAvailable = availableOpenIntervals
	nextBlocked = blockedOpenIntervals

	if e.IntervalType == Available {
		if e.EndpointType == Start {
			nextAvailable = availableOpenIntervals + 1
		} else {
			nextAvailable = availableOpenIntervals - 1
		}
	} else if e.IntervalType == Blocked {
		if e.EndpointType == Start {
			nextBlocked = blockedOpenIntervals + 1
		} else {
			nextBlocked = blockedOpenIntervals - 1
		}
	}
	return nextAvailable, nextBlocked
}

// if you have intervals that are overlapping, use this function to merge them. The resulting intervals will not overlap
func MergeAndReturnNonOverlappingIntervals(a []Interval) []Interval {
	return SubstractBlockedIntervals(a, []Interval{})
}

func IntervalsByDay(a []Interval) map[time.Time][]Interval {
	m := map[time.Time][]Interval{}

	for _, i := range a {
		if SameDay(i.Start, i.End) {
			addIntervalToMap(m, i)
		} else {
			addIntervalToMap(m, Interval{Start: i.Start, End: NormalizeDate(i.Start).AddDate(0, 0, 1).Add(-1 * time.Millisecond)})
			for cd := NormalizeDate(i.Start).AddDate(0, 0, 1); i.End.Sub(cd).Hours() >= 24; cd = cd.AddDate(0, 0, 1) {
				addIntervalToMap(m, Interval{Start: cd, End: cd.AddDate(0, 0, 1).Add(-1 * time.Millisecond)})
			}
			addIntervalToMap(m, Interval{Start: NormalizeDate(i.End), End: i.End})
		}
	}
	return m
}

func SameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

type DayIntervals struct {
	Date                     time.Time
	CountSinceFirst          int
	OrderedDisjunctIntervals []Interval
}

// Will return a list of non-overlapping DayIntervals ordered by DayIntervals.Date.
// Each Day will have a list of ordered disjoint intervals associated in that day
// An entry will be generated for each day in the range defined by the lowest start time and the highest start time covered by an interval
func IntervalsForEachDayInRange(a []Interval, startDay, endDay time.Time) ([]DayIntervals, error) {
	if startDay.Sub(endDay).Seconds() > 0 {
		return []DayIntervals{}, errors.New("Start day must be before the end")
	}
	if endDay.Sub(startDay).Hours() > 24*365 {
		return []DayIntervals{}, errors.New("Can not request more than 365 days")
	}

	m := IntervalsByDay(a)

	result := []DayIntervals{}
	d := 0
	for c := NormalizeDate(startDay); endDay.Sub(c) >= 0; c = c.AddDate(0, 0, 1) {
		result = append(result, DayIntervals{
			Date:                     c,
			CountSinceFirst:          d,
			OrderedDisjunctIntervals: m[c],
		})
		d ++
	}

	return result, nil
}

// returns a Time object with only the date component set
func NormalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func addIntervalToMap(m map[time.Time][]Interval, sameDayInterval Interval) {
	k := NormalizeDate(sameDayInterval.Start)
	m[k] = append(m[k], sameDayInterval)
}

func SplitInFixedIntervals(orderedDisjointIntervals []Interval, intervalLengthInMinutes int) []Interval {
	r := []Interval{}
	l := time.Duration(intervalLengthInMinutes)
	for _, i := range orderedDisjointIntervals {
		for c := i.Start; i.End.Add(-l*time.Minute).Sub(c).Seconds() > -1; c = c.Add(l*time.Minute) {
			r = append(r, Interval{Start: c, End: c.Add(l*time.Minute)})
		}
	}
	return r
}
