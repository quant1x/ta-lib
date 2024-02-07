package algorithms

import (
	"bytes"
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pandas/stat"
	"gitee.com/quant1x/ta-lib/plot"
	"github.com/wcharczuk/go-chart/v2" //exposes "chart"
	"log"
	"os"
	"testing"
)

func TestPeakDetect(t *testing.T) {
	data := []float64{1, 1, 1.1, 1, 0.9, 1, 1, 1.1, 1, 0.9, 1, 1.1, 1, 1, 0.9, 1, 1, 1.1, 1, 1, 1, 1, 1.1, 0.9, 1, 1.1, 1, 1, 0.9, 1, 1.1, 1, 1, 1.1, 1, 0.8, 0.9, 1, 1.2, 0.9, 1, 1, 1.1, 1.2, 1, 1.5, 1, 3, 2, 5, 3, 2, 1, 1, 1, 0.9, 1, 1, 3, 2.6, 4, 3, 3.2, 2, 1, 1, 0.8, 4, 4, 2, 2.5, 1, 1, 1}

	// Algorithm configuration from example.
	const (
		lag       = 30
		threshold = 5
		influence = 0
	)

	// Create then initialize the peak detector.
	detector := NewPeakDetector()
	err := detector.Initialize(influence, threshold, data[:lag]) // The length of the initial values is the lag.
	if err != nil {
		log.Fatalf("Failed to initialize peak detector.\nError: %s", err)
	}

	// Start processing new data points and determine what signal, if any they produce.
	//
	// This method, .Next(), is best for when data are being processed in a stream, but this simply iterates over a
	// slice.
	nextDataPoints := data[lag:]
	for i, newPoint := range nextDataPoints {
		signal := detector.Next(newPoint)
		var signalType string
		switch signal {
		case SignalNegative:
			signalType = "negative"
		case SignalNeutral:
			signalType = "neutral"
		case SignalPositive:
			signalType = "positive"
		}

		println(fmt.Sprintf("Data point at index %d has the signal: %s", i+lag, signalType))
	}

	// This method, .NextBatch(), is a helper function for processing many data points at once. It's returned slice
	// should produce the same signal outputs as the loop above.
	signals := detector.NextBatch(nextDataPoints)
	println(fmt.Sprintf("1:1 ratio of batch inputs to signal outputs: %t", len(signals) == len(nextDataPoints)))
}

func TestPeek(t *testing.T) {
	const (
		lag       = 30
		threshold = 5
		influence = 0
	)
	code := "sh000001"
	date := "20240206"
	klines := base.CheckoutKLines(code, date)
	if len(klines) == 0 {
		return
	}
	df := pandas.LoadStructs(klines)
	if df.Nrow() == 0 {
		return
	}
	fmt.Println(df)
	CLOSE := df.ColAsNDArray("close")
	data := CLOSE.DTypes()
	detector := NewPeakDetector()
	err := detector.Initialize(influence, threshold, data) // The length of the initial values is the lag.
	if err != nil {
		log.Fatalf("Failed to initialize peak detector.\nError: %s", err)
	}
	signals := detector.NextBatch(data)
	println(fmt.Sprintf("1:1 ratio of batch inputs to signal outputs: %t", len(signals) == len(data)))
	// 输出图表
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: data,
				YValues: data,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.PNG, buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data)
	err = os.WriteFile(code+".png", buffer.Bytes(), api.CACHE_FILE_MODE)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestPeekForGoCharts(t *testing.T) {
	const (
		lag       = 30
		threshold = 0
		influence = 0
	)
	code := "sh000001"
	date := "20240206"
	klines := base.CheckoutKLines(code, date)
	if len(klines) == 0 {
		return
	}
	df := pandas.LoadStructs(klines)
	if df.Nrow() == 0 {
		return
	}
	rows := 89
	sl := api.RangeFinite(df.Nrow() - rows)
	df = df.SelectRows(sl)
	fmt.Println(df)
	DATE := df.Col("date").Strings()
	CLOSE := df.ColAsNDArray("close")
	//HIGH := df.ColAsNDArray("high")
	data := CLOSE.DTypes()
	detector := NewPeakDetector()
	err := detector.Initialize(influence, threshold, data) // The length of the initial values is the lag.
	if err != nil {
		log.Fatalf("Failed to initialize peak detector.\nError: %s", err)
	}
	signals := detector.NextBatch(data)
	fmt.Println(signals)
	println(fmt.Sprintf("1:1 ratio of batch inputs to signal outputs: %t", len(signals) == len(data)))
	// 输出图表
	x := stat.Range[float64](CLOSE.Len())
	font, _ := plot.GetDefaultFont()
	peeks := make([]float64, 0, rows)
	frontValue := float64(0)
	for i, v := range signals {
		f := data[i]
		//if frontValue != 0 {
		//	f = frontValue
		//}
		//
		//if v != SignalNeutral {
		//	frontValue = f
		//}
		peeks = append(peeks, f)
		_ = v
	}
	_ = frontValue
	graph := chart.Chart{
		Title: "走势图",
		Font:  font,
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "收盘价",
				XValues: x,
				YValues: CLOSE.DTypes(),
				XValueFormatter: func(v interface{}) string {
					f := v.(float64)
					idx := int(f)
					return DATE[idx]
				},
			},
			chart.ContinuousSeries{
				XValues: x,
				YValues: peeks,
				XValueFormatter: func(v interface{}) string {
					f := v.(float64)
					idx := int(f)
					return DATE[idx]
				},
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.PNG, buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data)
	err = os.WriteFile(code+".png", buffer.Bytes(), api.CACHE_FILE_MODE)
	if err != nil {
		fmt.Println(err)
		return
	}
}
