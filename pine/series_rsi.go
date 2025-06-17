package pine

// RSI generates a ValueSeries of relative strength index
//
// The formula for RSI is
//   - u = Count the number of p(t+1) - p(t) > 0 as gains
//   - d = Count the number of p(t+1) - p(t) < 0 as losses
//   - rs = ta.rma(u) / ta.rma(d)
//   - res = 100 - 100 / (1 + rs)
func RSI(p ValueSeries, l int64) ValueSeries {
	changes := Change(p, 1)

	upper := OperateNoCache(changes, changes, "upper", func(a, b float64) float64 {
		return max(a, 0)
	})

	lower := OperateNoCache(changes, changes, "lower", func(a, b float64) float64 {
		return max(-a, 0)
	})

	upperRMA := RMA(upper, l)
	lowerRMA := RMA(lower, l)

	rsi := OperateNoCache(upperRMA, lowerRMA, "rsi", func(a, b float64) float64 {
		if b == 0 {
			return 100
		}

		return 100 - 100/(1+a/b)
	})

	return rsi
}
