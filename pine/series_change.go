package pine

// Change compares the current `source` value to its value `lookback` bars ago and returns the difference.
//
// arguments are
//   - src: ValueSeries - Source data to seek difference
//   - lookback: int - Lookback to compare the change
func Change(src ValueSeries, lookback int) ValueSeries {
	change := NewValueSeries()
	stop := src.GetCurrent()
	if stop == nil {
		return change
	}

	value := src.GetFirst()
	for value != nil {
		if changeValue := getChangeValue(value, lookback); changeValue != nil {
			change.Set(value.t, *changeValue)
		}

		if value.t.Equal(stop.t) {
			break
		}

		value = value.next
	}
	change.SetCurrent(stop.t)
	return change
}

func getChangeValue(value *Value, lookback int) *float64 {
	current := value.v
	for range lookback {
		value = value.prev
		if value == nil {
			return nil
		}
	}

	return NewFloat64(current - value.v)
}

func NewFloat64(v float64) *float64 {
	v2 := v
	return &v2
}
