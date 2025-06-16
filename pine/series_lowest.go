package pine

// Lowest generates a ValueSeries of lowest value of previous values
//
// Parameters
//   - p - ValueSeries: source data
//   - l - int: lookback periods [1, ∞)
func Lowest(p ValueSeries, l int) ValueSeries {
	lowest := NewValueSeries()
	stop := p.GetCurrent()
	if stop == nil {
		return lowest
	}

	value := p.GetFirst()
	for value != nil {
		if lowestValue := getLowestValue(value, l); lowestValue != nil {
			lowest.Set(value.t, *lowestValue)
		}

		if value.t.Equal(stop.t) {
			break
		}

		value = value.next
	}

	lowest.SetCurrent(stop.t)
	return lowest
}

func getLowestValue(value *Value, l int) *float64 {
	var lowest *float64
	for i := 0; i < l; i++ {
		value = value.prev
		if value == nil {
			return nil
		}

		if lowest == nil || value.v < *lowest {
			lowest = &value.v
		}
	}

	return lowest
}
