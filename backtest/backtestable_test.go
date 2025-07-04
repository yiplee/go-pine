/*
Backtesting is a process of evaluating a trading strategy using historical market data to simulate how the strategy would have performed in the past.
This library provides such capability using the PineScript-like indicators.
*/
package backtest

import (
	"log"
	"time"

	"github.com/yiplee/go-pine/pine"
)

type mystrat struct{}

func (m *mystrat) OnNextOHLCV(strategy Strategy, s pine.OHLCVSeries, state map[string]interface{}) error {

	close := pine.OHLCVAttr(s, pine.OHLCPropClose)
	rsi := pine.SMA(close, 2)
	macd, _, _ := pine.MACD(close, 12, 26, 9)
	stdev := pine.Stdev(close, 24)
	ema200 := pine.EMA(close, 200)

	// we haven't seen enough candles to fulfill the lookback period
	if rsi.Val() == nil || macd.Val() == nil || stdev.Val() == nil || ema200.Val() == nil {
		return nil
	}

	if *rsi.Val() < 30 && *macd.Val() < 0 && *ema200.Val() > 0 {
		entry1 := EntryOpts{
			Side: Long,
		}
		strategy.Entry("Buy1", entry1)
	}

	if *rsi.Val() > 70 && *macd.Val() > 0 {
		strategy.Exit("Buy1")
	}

	return nil
}

func ExampleBackTestable() {
	b := &mystrat{}
	data := pine.OHLCVTestData(time.Now(), 25, 5*60*1000)
	series, _ := pine.NewOHLCVSeries(data)

	res, _ := RunBacktest(series, b)

	log.Printf("TotalClosedTrades %d, PercentProfitable: %.03f, NetProfit: %s", res.TotalClosedTrades, res.PercentProfitable, res.NetProfit)

}
