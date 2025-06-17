package pine

import (
	"log"
	"testing"
	"time"
)

// TestSeriesHighestNoData tests no data scenario
func TestSeriesHighestNoData(t *testing.T) {
	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	prop := OHLCVAttr(series, OHLCPropClose)
	highest := Highest(prop, 2)
	if highest == nil {
		t.Error("Expected to be non nil but got nil")
	}
}

// TestSeriesHighestNoIteration tests scenario where there's no iteration yet
func TestSeriesHighestNoIteration(t *testing.T) {
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
	highest := Highest(prop, 2)
	if highest == nil {
		t.Error("Expected to be non-nil but got nil")
	}
}

// TestSeriesHighestIteration tests iteration scenario
func TestSeriesHighestIteration(t *testing.T) {
	start := time.Now()
	data := OHLCVTestData(start, 5, 5*60*1000)
	data[0].C = 13
	data[1].C = 15
	data[2].C = 11
	data[3].C = 19
	data[4].C = 21

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		lookback int
		vals     []float64
	}{
		{
			lookback: 1,
			vals:     []float64{13, 15, 11, 19, 21},
		},
		{
			lookback: 2,
			vals:     []float64{0, 15, 15, 19, 21},
		},
		{
			lookback: 3,
			vals:     []float64{0, 0, 15, 19, 21},
		},
	}

	for j := 0; j <= 3; j++ {
		series.Next()

		for i, v := range testTable {
			prop := OHLCVAttr(series, OHLCPropClose)
			highest := Highest(prop, v.lookback)
			exp := v.vals[j]
			if exp == 0 {
				if highest.Val() != nil {
					t.Fatalf("expected nil but got non nil: %+v at vals item: %d, testtable item: %d", *highest.Val(), j, i)
				}
				// OK
			}
			if exp != 0 {
				if highest.Val() == nil {
					t.Fatalf("expected non nil: %+v but got nil at vals item: %d, testtable item: %d", exp, j, i)
				}
				if exp != *highest.Val() {
					t.Fatalf("expected %+v but got %+v at vals item: %d, testtable item: %d", exp, *highest.Val(), j, i)
				}
				// OK
			}
		}
	}
}

func TestMemoryLeakHighest(t *testing.T) {
	testMemoryLeak(t, func(o OHLCVSeries) error {
		prop := OHLCVAttr(o, OHLCPropClose)
		Highest(prop, 10)
		return nil
	})
}

func ExampleHighest() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	for {
		if v, _ := series.Next(); v == nil {
			break
		}

		close := OHLCVAttr(series, OHLCPropClose)
		// Get the highest of last 10 values
		highest := Highest(close, 10)
		log.Printf("Highest: %+v", highest.Val())
	}
}
