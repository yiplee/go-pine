package pine

import (
	"log"
	"testing"
	"time"
)

// TestSeriesRSINoData tests no data scenario
//
// t=time.Time (no iteration) | |
// p=ValueSeries              | |
// rsi=ValueSeries            | |
func TestSeriesRSINoData(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	prop := OHLCVAttr(series, OHLCPropClose)
	rsi := RSI(prop, 2)
	if rsi == nil {
		t.Error("Expected to be non nil but got nil")
	}
}

// TestSeriesRSINoIteration tests this sceneario where there's no iteration yet
//
// t=time.Time (no iteration) | 1  |  2   | 3  | 4  |
// p=ValueSeries              | 14 |  15  | 17 | 18 |
// rsi=ValueSeries            |    |      |    |    |
func TestSeriesRSINoIteration(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)
	data[0].C = 14
	data[1].C = 15
	data[2].C = 17
	data[3].C = 18

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	prop := OHLCVAttr(series, OHLCPropClose)
	rsi := RSI(prop, 2)
	if rsi == nil {
		t.Error("Expected to be non-nil but got nil")
	}
}

// TestSeriesRSIIteration5 tests this scneario when the iterator is at t=4 is not at the end
//
// t=time.Time (no iteration) | 1   |  2  | 3    | 4 (here1)  | 5 (here2)  |
// p=ValueSeries              | 13  | 15  | 11   | 18         | 20        |
// u(close, 2)                | nil | nil |  2   | 7          | 9         |
// d(close, 2)                | nil | nil |  4   | 4          | 0         |
// rsi(u(close,2), 2)		  | nil | nil | nil  | 4.5        | 6.75      |
// rsi(d(close,2), 2)		  | nil | nil | nil  | 4          | 2         |
// rsi(close, 2)			  | nil | nil | nil| | 52.9411765 | 77.1428571|
func TestSeriesRSIIteration5(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 5, 5*60*1000)
	data[0].C = 13
	data[1].C = 15
	data[2].C = 11
	data[3].C = 18
	data[4].C = 20

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	series.Next()
	series.Next()
	series.Next()

	testTable := []float64{80, 85.71428571428571}

	for i, v := range testTable {
		series.Next()

		prop := OHLCVAttr(series, OHLCPropClose)
		rsi := RSI(prop, 2)
		if rsi == nil {
			t.Errorf("Expected to be non nil but got nil at idx: %d", i)
		}
		if *rsi.Val() != v {
			t.Errorf("Expected %+v but got %+v at idx: %d", v, *rsi.Val(), i)
		}
	}
}

// TestSeriesRSINotEnoughData tests this scneario when the lookback is more than the number of data available
//
// t=time.Time    | 1  |  2   | 3  | 4 (here)  |
// p=ValueSeries  | 14 |  15  | 17 | 18        |
// rsi(close, 5)  | nil| nil  | nil| nil       |
func TestSeriesRSINotEnoughData(t *testing.T) {

	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)
	data[0].C = 13
	data[1].C = 15
	data[2].C = 11
	data[3].C = 18

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	series.Next()
	series.Next()
	series.Next()
	series.Next()

	testTable := []struct {
		lookback int
		exp      *float64
	}{
		{
			lookback: 5,
			exp:      nil,
		},
		{
			lookback: 6,
			exp:      nil,
		},
	}

	for i, v := range testTable {
		prop := OHLCVAttr(series, OHLCPropClose)

		rsi := RSI(prop, int64(v.lookback))
		if rsi == nil {
			t.Errorf("Expected to be non nil but got nil at idx: %d", i)
		}
		if rsi.Val() != v.exp {
			t.Errorf("Expected to get %+v but got %+v for lookback %+v", v.exp, *rsi.Val(), v.lookback)
		}
	}
}

func TestMemoryLeakRSI(t *testing.T) {
	testMemoryLeak(t, func(o OHLCVSeries) error {
		prop := OHLCVAttr(o, OHLCPropClose)
		RSI(prop, 12)
		return nil
	})
}

func ExampleRSI() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	for {
		if v, _ := series.Next(); v == nil {
			break
		}

		close := OHLCVAttr(series, OHLCPropClose)
		rsi := RSI(close, 16)
		log.Printf("RSI: %+v", rsi.Val())
	}
}
