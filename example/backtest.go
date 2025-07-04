package main

import (
	"log"
	"time"

	"github.com/yiplee/go-pine/backtest"
	"github.com/yiplee/go-pine/pine"
)

type mystrat struct{}

func (m *mystrat) OnNextOHLCV(strategy backtest.Strategy, s pine.OHLCVSeries, state map[string]interface{}) error {

	var short int64 = 2
	var long int64 = 20
	var span float64 = 10

	close := pine.OHLCVAttr(s, pine.OHLCPropClose)
	open := pine.OHLCVAttr(s, pine.OHLCPropOpen)

	basis := pine.SMA(close, short)
	basis2 := pine.SMA(open, long)
	rsi := pine.RSI(close, short)
	avg := pine.SMA(rsi, long)

	basis3 := pine.Add(basis, basis2)
	upperBB := pine.AddConst(basis3, span)

	log.Printf("t: %+v, close: %+v, rsi: %+v, avg: %+v, upperBB: %+v", s.Current().S, close.Val(), rsi.Val(), avg.Val(), upperBB.Val())

	if rsi.Val() != nil {
		if *rsi.Val() < 30 {
			log.Printf("Entry: %+v", *rsi.Val())
			entry1 := backtest.EntryOpts{
				Side: backtest.Long,
			}
			strategy.Entry("Buy1", entry1)
		}
		if *rsi.Val() > 70 {
			log.Printf("Exit %+v", *rsi.Val())
			strategy.Exit("Buy1")
		}
	}

	return nil
}

func main() {
	b := &mystrat{}
	data := pine.OHLCVTestData(time.Now(), 25, 5*60*1000)
	series, _ := pine.NewOHLCVSeries(data)

	res, _ := backtest.RunBacktest(series, b)
	log.Printf("TotalClosedTrades %d, PercentProfitable: %.03f, NetProfit: %s", res.TotalClosedTrades, res.PercentProfitable, res.NetProfit)
}
