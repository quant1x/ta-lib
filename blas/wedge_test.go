package blas

import (
	"bytes"
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/ta-lib/plot"
	"github.com/wcharczuk/go-chart/v2"
	"os"
	"testing"
)

func TestWedge_basic(t *testing.T) {
	requiredKLines := 89
	requiredKLines = 34
	//requiredKLines = 250
	code := "sh000001"
	//code = "sz300629"
	//code = "000917"
	//code = "600843"
	date := "2024-04-01"
	//date = "2024-03-29"
	//date = cache.DefaultCanReadDate()
	list := base.CheckoutKLines(code, date)
	if len(list) >= requiredKLines {
		list = list[len(list)-requiredKLines:]
	}
	sample := LoadKLineSample(list)
	rows := sample.Len()
	securityCode := exchange.CorrectSecurityCode(code)
	waves := NewWaves(sample, securityCode)
	fmt.Println(waves)
	ticks := make([]chart.Tick, rows)
	data := make([]float64, rows)
	volumes := make([]float64, rows)
	DATE := make([]string, rows)
	for i := 0; i < rows; i++ {
		v := list[i]
		DATE[i] = v.Date
		ticks[i] = chart.Tick{Value: float64(i), Label: DATE[i]}
		data[i] = v.Close
		volumes[i] = v.Volume
		//if i > 0 {
		//	if list[i-1].High < v.High {
		//		data[i] = v.High
		//	} else if list[i-1].Low > v.Low {
		//		data[i] = v.Low
		//	}
		//}
		//if i > 0 {
		//	if list[i-1].Close < v.Close && list[i-1].High < v.High {
		//		data[i] = v.High
		//	} else if list[i-1].Close > v.Close && list[i-1].Low > v.Low {
		//		data[i] = v.Low
		//	}
		//}
	}
	xAxisFormat := func(v interface{}) string {
		f := v.(float64)
		idx := int(f)
		x := DATE[idx]
		return x
	}
	graph := plot.CreateChart()
	graph.XAxis = chart.XAxis{
		Name:           securities.GetStockName(code) + "(" + securityCode + ")日线图 - " + date,
		ValueFormatter: xAxisFormat,
		TickStyle: chart.Style{
			TextRotationDegrees: 45.0,
		},
		Ticks: ticks,
	}
	graph.YAxis = chart.YAxis{
		Name: "股价",
		GridMajorStyle: chart.Style{
			StrokeColor:     chart.ColorAlternateLightGray.WithAlpha(200),
			StrokeWidth:     1.0,
			StrokeDashArray: plot.DashedLine,
		},
	}
	dataY := data
	x := num.Range[float64](rows)
	dataX := x
	graph.Series = []chart.Series{
		chart.ContinuousSeries{
			Name:            "CLOSE",
			XValues:         dataX,
			YValues:         dataY,
			XValueFormatter: xAxisFormat,
			Style:           chart.Style{StrokeColor: chart.ColorBlack},
		},
		chart.ContinuousSeries{
			Name:            "VOLUME",
			YAxis:           chart.YAxisSecondary,
			XValues:         dataX,
			YValues:         volumes,
			XValueFormatter: xAxisFormat,
			Style:           chart.Style{StrokeColor: chart.ColorAlternateGray, StrokeDashArray: plot.DashedLine},
		},
	}

	pattern := MatchWedge(waves)
	if pattern != nil {
		fmt.Printf("wedge=%+v\n", pattern)
		series := pattern.ExportSeries(sample)
		for i := 0; i < len(series); i++ {
			//serieses[i].XValueFormatter = xAxisFormat
			graph.Series = append(graph.Series, series[i])
		}
	}
	graph.Elements = []chart.Renderable{chart.LegendThin(&graph)}
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data)
	err = os.WriteFile("wedge-kline-"+securityCode+".png", buffer.Bytes(), api.CACHE_FILE_MODE)
	if err != nil {
		fmt.Println(err)
		return
	}
}
