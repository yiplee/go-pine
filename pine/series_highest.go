package pine

import (
	"fmt"
)

// Highest generates a ValueSeries of highest value of previous values
//
// Parameters
//   - p - ValueSeries: source data
//   - l - int: lookback periods [1, ∞)
func Highest(p ValueSeries, l int) ValueSeries {

	key := fmt.Sprintf("highest:%s:%d", p.ID(), l)
	highest := getCache(key)
	if highest == nil {
		highest = NewValueSeries()
	}

	highest = generateHighest(p, highest, l)

	setCache(key, highest)

	return highest
}

// HighestNoCache generates highest without caching
func HighestNoCache(p ValueSeries, l int) ValueSeries {
	highest := NewValueSeries()
	return generateHighest(p, highest, l)
}

func generateHighest(p, highest ValueSeries, l int) ValueSeries {
	// current available value
	stop := p.GetCurrent()
	if stop == nil {
		return highest
	}
	highest = getHighest(*stop, highest, p, l)
	highest.SetCurrent(stop.t)
	return highest
}

func getHighest(stop Value, highest ValueSeries, src ValueSeries, l int) ValueSeries {
	// keep track of the source values of highest, maximum of l+1 items
	highestSrc := make([]float64, 0)
	var startNew *Value

	lastAvail := highest.GetLast()

	if lastAvail == nil {
		startNew = src.GetFirst()
	} else {
		v := src.Get(lastAvail.t)
		startNew = v.next
	}

	if startNew == nil {
		// if nothing is to start with, then nothing can be done
		return highest
	}

	// populate source values to be checked for highest
	if lastAvail != nil {
		lastAvailv := src.Get(lastAvail.t)

		for {
			if lastAvailv == nil {
				break
			}

			srcv := src.Get(lastAvailv.t)
			// add at the beginning since we go backwards
			highestSrc = append([]float64{srcv.v}, highestSrc...)

			if len(highestSrc) == l {
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
		highestSrc = append(highestSrc, v.v)

		var set bool

		// if previous exists, we just update the window and set the new highest
		if v.prev != nil {
			e := highest.Get(v.prev.t)
			if e != nil && len(highestSrc) == l+1 {
				// remove the oldest value and find the max in the window
				window := highestSrc[1:]
				maxVal := window[0]
				for _, val := range window {
					if val > maxVal {
						maxVal = val
					}
				}
				highest.Set(itervt, maxVal)
				set = true
			}
		}

		if !set {
			if len(highestSrc) >= l {
				var ct int
				var maxVal float64
				for i := len(highestSrc) - 1; i >= 0; i-- {
					if ct == 0 {
						maxVal = highestSrc[i]
					} else if highestSrc[i] > maxVal {
						maxVal = highestSrc[i]
					}
					ct++
					if ct == l {
						break
					}
				}
				highest.Set(itervt, maxVal)
			}
		}

		if v.next == nil {
			break
		}
		if v.t.Equal(stop.t) {
			break
		}

		if len(highestSrc) > l+1 {
			highestSrc = highestSrc[1:]
		}
		itervt = v.next.t
	}

	return highest
}
