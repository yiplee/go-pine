package pine

import (
	"fmt"
)

// Lowest generates a ValueSeries of lowest value of previous values
//
// Parameters
//   - p - ValueSeries: source data
//   - l - int: lookback periods [1, ∞)
func Lowest(p ValueSeries, l int) ValueSeries {
	key := fmt.Sprintf("lowest:%s:%d", p.ID(), l)
	lowest := getCache(key)
	if lowest == nil {
		lowest = NewValueSeries()
	}

	lowest = generateLowest(p, lowest, l)

	setCache(key, lowest)

	return lowest
}

// LowestNoCache generates lowest without caching
func LowestNoCache(p ValueSeries, l int) ValueSeries {
	lowest := NewValueSeries()
	return generateLowest(p, lowest, l)
}

func generateLowest(p, lowest ValueSeries, l int) ValueSeries {
	// current available value
	stop := p.GetCurrent()
	if stop == nil {
		return lowest
	}
	lowest = getLowest(*stop, lowest, p, l)
	lowest.SetCurrent(stop.t)
	return lowest
}

func getLowest(stop Value, lowest ValueSeries, src ValueSeries, l int) ValueSeries {
	// keep track of the source values of lowest, maximum of l+1 items
	lowestSrc := make([]float64, 0)
	var startNew *Value

	lastAvail := lowest.GetLast()

	if lastAvail == nil {
		startNew = src.GetFirst()
	} else {
		v := src.Get(lastAvail.t)
		startNew = v.next
	}

	if startNew == nil {
		// if nothing is to start with, then nothing can be done
		return lowest
	}

	// populate source values to be checked for lowest
	if lastAvail != nil {
		lastAvailv := src.Get(lastAvail.t)

		for {
			if lastAvailv == nil {
				break
			}

			srcv := src.Get(lastAvailv.t)
			// add at the beginning since we go backwards
			lowestSrc = append([]float64{srcv.v}, lowestSrc...)

			if len(lowestSrc) == l {
				break
			}
			lastAvailv = lastAvailv.prev
		}
	}

	// first new time
	itervt := startNew.t

	for {
		v := src.Get(itervt)
		if v == nil {
			break
		}

		// append new source data
		lowestSrc = append(lowestSrc, v.v)

		var set bool

		// if previous exists, we just update the window and set the new lowest
		if v.prev != nil {
			e := lowest.Get(v.prev.t)
			if e != nil && len(lowestSrc) == l+1 {
				// remove the oldest value and find the min in the window
				window := lowestSrc[1:]
				minVal := window[0]
				for _, val := range window {
					if val < minVal {
						minVal = val
					}
				}
				lowest.Set(itervt, minVal)
				set = true
			}
		}

		if !set {
			if len(lowestSrc) >= l {
				var ct int
				var minVal float64
				for i := len(lowestSrc) - 1; i >= 0; i-- {
					if ct == 0 {
						minVal = lowestSrc[i]
					} else if lowestSrc[i] < minVal {
						minVal = lowestSrc[i]
					}
					ct++
					if ct == l {
						break
					}
				}
				lowest.Set(itervt, minVal)
			}
		}

		if v.next == nil {
			break
		}
		if v.t.Equal(stop.t) {
			break
		}

		if len(lowestSrc) > l+1 {
			lowestSrc = lowestSrc[1:]
		}
		itervt = v.next.t
	}

	return lowest
}
