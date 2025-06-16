package pine

// Highest generates a ValueSeries of highest value of previous values
//
// Parameters
//   - p - ValueSeries: source data
//   - l - int: lookback periods [1, ∞)
func Highest(p ValueSeries, l int) ValueSeries {
	highest := NewValueSeries()
	stop := p.GetCurrent()
	if stop == nil {
		return highest
	}

	value := p.GetFirst()
	for value != nil {
		if highestValue := getHighestValue(value, l); highestValue != nil {
			highest.Set(value.t, *highestValue)
		}

		if value.t.Equal(stop.t) {
			break
		}

		value = value.next
	}

	highest.SetCurrent(stop.t)
	return highest
}

// get the highest value of the previous l values
func getHighestValue(value *Value, l int) *float64 {
	var highest *float64
	for i := 0; i < l; i++ {
		value = value.prev
		if value == nil {
			return nil
		}

		if highest == nil || value.v > *highest {
			highest = &value.v
		}
	}
	return highest
}
