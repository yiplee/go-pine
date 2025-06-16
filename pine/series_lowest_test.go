package pine

import (
	"log"
	"testing"
	"time"
)

// TestSeriesLowestNoData tests no data scenario
func TestSeriesLowestNoData(t *testing.T) {
	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	prop := OHLCVAttr(series, OHLCPropLow)
	lowest := Lowest(prop, 2)
	if lowest == nil {
		t.Error("Expected to be non nil but got nil")
	}
}

// TestSeriesLowestNoIteration tests scenario where there's no iteration yet
func TestSeriesLowestNoIteration(t *testing.T) {
	start := time.Now()
	data := OHLCVTestData(start, 4, 5*60*1000)
	data[0].L = 14
	data[1].L = 15
	data[2].L = 17
	data[3].L = 18

	series, err := NewOHLCVSeries(data)
	if err != nil {
		t.Fatal(err)
	}

	prop := OHLCVAttr(series, OHLCPropLow)
	lowest := Lowest(prop, 2)
	if lowest == nil {
		t.Error("Expected to be non-nil but got nil")
	}
}

// TestSeriesLowestIteration tests iteration scenario
func TestSeriesLowestIteration(t *testing.T) {
	start := time.Now()
	data := OHLCVTestData(start, 5, 5*60*1000)
	data[0].L = 13
	data[1].L = 15
	data[2].L = 11
	data[3].L = 19
	data[4].L = 21

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
			vals:     []float64{0, 13, 15, 11, 19},
		},
		{
			lookback: 2,
			vals:     []float64{0, 0, 13, 11, 11},
		},
		{
			lookback: 3,
			vals:     []float64{0, 0, 0, 11, 11},
		},
	}

	for j := 0; j <= 3; j++ {
		series.Next()

		for i, v := range testTable {
			prop := OHLCVAttr(series, OHLCPropLow)
			lowest := Lowest(prop, v.lookback)
			exp := v.vals[j]
			if exp == 0 {
				if lowest.Val() != nil {
					t.Fatalf("expected nil but got non nil: %+v at vals item: %d, testtable item: %d", *lowest.Val(), j, i)
				}
				// OK
			}
			if exp != 0 {
				if lowest.Val() == nil {
					t.Fatalf("expected non nil: %+v but got nil at vals item: %d, testtable item: %d", exp, j, i)
				}
				if exp != *lowest.Val() {
					t.Fatalf("expected %+v but got %+v at vals item: %d, testtable item: %d", exp, *lowest.Val(), j, i)
				}
				// OK
			}
		}
	}
}

func TestMemoryLeakLowest(t *testing.T) {
	testMemoryLeak(t, func(o OHLCVSeries) error {
		prop := OHLCVAttr(o, OHLCPropLow)
		Lowest(prop, 10)
		return nil
	})
}

func ExampleLowest() {
	start := time.Now()
	data := OHLCVTestData(start, 10000, 5*60*1000)
	series, _ := NewOHLCVSeries(data)
	for {
		if v, _ := series.Next(); v == nil {
			break
		}

		low := OHLCVAttr(series, OHLCPropLow)
		// Get the lowest of last 10 values
		lowest := Lowest(low, 10)
		log.Printf("Lowest: %+v", lowest.Val())
	}
}
