package pine

type smaCalcItem struct {
	valuetot float64
	total    int64
	seeked   int64
}

// SMA generates a ValueSeries of simple moving averages
func SMA(p ValueSeries, l int64) ValueSeries {
	sma := NewValueSeries()
	stop := p.GetCurrent()
	if stop == nil {
		return sma
	}

	value := p.GetFirst()
	for value != nil {
		if smaValue := getSMAValue(value, l); smaValue != nil {
			sma.Set(value.t, *smaValue)
		}

		if value.t.Equal(stop.t) {
			break
		}

		value = value.next
	}
	sma.SetCurrent(stop.t)
	return sma
}

func getSMAValue(value *Value, l int64) *float64 {
	sum := value.v
	for range l - 1 {
		value = value.prev
		if value == nil {
			return nil
		}

		sum += value.v
	}
	return NewFloat64(sum / float64(l))
}
