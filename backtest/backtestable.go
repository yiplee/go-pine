package backtest

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/yiplee/go-pine/pine"
)

type BackTestable interface {
	OnNextOHLCV(Strategy, pine.OHLCVSeries, map[string]interface{}) error
}

type BacktestResult struct {
	ClosedOrd         []Position
	NetProfit         decimal.Decimal
	PercentProfitable float64
	ProfitableTrades  int64
	TotalClosedTrades int64
}

// EntryOpts is additional entry options
type EntryOpts struct {
	Comment string

	// Limit price is used if this value is non nil. If it's nil, market order is executed
	Limit *float64

	OrdID string
	Side  Side
	Stop  string
	Qty   string
}

// Px generates a non nil float64
func Px(v float64) *float64 {
	v2 := &v
	return v2
}

type Side string

const (
	Long  Side = "long"
	Short Side = "short"
)

type Position struct {
	Qty       decimal.Decimal
	EntryPx   decimal.Decimal
	ExitPx    decimal.Decimal
	EntryTime time.Time
	ExitTime  time.Time
	EntrySide Side
	OrdID     string
}

func (p Position) Profit() decimal.Decimal {
	switch p.EntrySide {
	case Long:
		return p.ExitPx.Sub(p.EntryPx).Mul(p.Qty)
	case Short:
		return p.EntryPx.Sub(p.ExitPx).Mul(p.Qty)
	}
	return decimal.Zero
}

func (b *BacktestResult) CalculateNetProfit() {
	start := decimal.NewFromFloat(0)
	for _, v := range b.ClosedOrd {
		p := v.Profit()
		start = start.Add(p)
	}
	b.NetProfit = start
}
