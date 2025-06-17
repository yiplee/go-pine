package pine

// RMA generates a ValueSeries of the exponentially weighted moving average with alpha = 1 / length.
// This is equivalent to J. Welles Wilder's smoothed moving average.
func RMA(p ValueSeries, l int64) ValueSeries {
	rma := NewValueSeries()
	stop := p.GetCurrent()
	if stop == nil {
		return rma
	}

	alpha := 1 / float64(l)
	value := p.GetFirst()
	for value != nil {
		if last := rma.GetLast(); last != nil {
			rma.Set(value.t, alpha*value.v+(1-alpha)*last.v)
		} else if sma := getSMAValue(value, l); sma != nil {
			rma.Set(value.t, *sma)
		}

		if value.t.Equal(stop.t) {
			break
		}

		value = value.next
	}
	rma.SetCurrent(stop.t)
	return rma
}
