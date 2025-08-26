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
	lowest := value.v
	for range l - 1 {
		value = value.prev
		if value == nil {
			return nil
		}

		if value.v < lowest {
			lowest = value.v
		}
	}

	return &lowest
}

func GetLowestValue(p ValueSeries, l int) *float64 {
	cur := p.GetCurrent()
	if cur == nil {
		return nil
	}

	return getLowestValue(cur, l)
}
