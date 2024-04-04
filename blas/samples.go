package blas

import (
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/ta-lib/plot"
	"strings"
)

// DataSample 数据样本
type DataSample interface {
	Time(n int) string             // 时间
	Len() int                      // 数据长度
	Current(n int) float64         // 当前值
	High(n int) float64            // 最高
	Low(n int) float64             // 最低
	Volume(n int) float64          // 量
	Chart(name string) *plot.Chart // 图表
}

// KLineSample K线样板的实现
type KLineSample struct {
	data []base.KLine
}

// LoadKLineSample 加载K线样本
func LoadKLineSample(data []base.KLine) KLineSample {
	return KLineSample{data: data}
}

func (k KLineSample) Time(n int) string {
	return k.data[n].Date
}

func (k KLineSample) Len() int {
	return len(k.data)
}

func (k KLineSample) value(n int, fieldName string) float64 {
	fieldName = strings.ToLower(fieldName)
	switch fieldName {
	case "open":
		return k.data[n].Open
	case "close":
		return k.data[n].Close
	case "high":
		return k.data[n].High
	case "low":
		return k.data[n].Low
	}
	return k.Current(n)
}

func (k KLineSample) Current(n int) float64 {
	return k.data[n].Close
}

func (k KLineSample) High(n int) float64 {
	wave := config.GetDataConfig().Feature.Wave
	return k.value(n, wave.Fields.Peak)
}

func (k KLineSample) Low(n int) float64 {
	wave := config.GetDataConfig().Feature.Wave
	return k.value(n, wave.Fields.Valley)
}

func (k KLineSample) Volume(n int) float64 {
	return k.data[n].Volume
}

func (k KLineSample) Chart(name string) *plot.Chart {
	xAxisFormat := func(v any) string {
		f := v.(float64)
		idx := int(f)
		//xAxis := DATE[idx]
		x := k.Time(idx)
		return x
	}
	rows := k.Len()
	xAxis := num.Range[float64](rows)
	yPrimaryName := "PRICE"
	yPrimary := make([]float64, rows)
	ySecondaryName := "VOLUME"
	ySecondary := make([]float64, rows)
	//ticks := make([]plot.Tick, plot.XTickMax)
	ticks := make([]plot.Tick, rows)
	for i := 0; i < rows; i++ {
		yPrimary[i] = k.Current(i)
		ySecondary[i] = k.Volume(i)
		//ti := int(float64(i) / float64(rows) * plot.XTickMax)
		//ticks[ti] = plot.Tick{Value: float64(i), Label: k.Time(i)}
		ticks[i] = plot.Tick{Value: float64(i), Label: k.Time(i)}
	}
	graph := plot.NewChart()
	graph.XAxis = plot.XAxis{
		Name:           name,
		TickPosition:   plot.TickPositionUnset,
		ValueFormatter: xAxisFormat,
		TickStyle: plot.Style{
			TextRotationDegrees: 45.0,
		},
		Ticks: ticks,
	}
	graph.YAxis = plot.YAxis{
		Name: "股价",
		GridMajorStyle: plot.Style{
			StrokeColor:     plot.ColorAlternateLightGray.WithAlpha(200),
			StrokeWidth:     1.0,
			StrokeDashArray: plot.DashedLine,
		},
	}
	graph.Series = []plot.Series{
		plot.ContinuousSeries{
			Name:    yPrimaryName,
			XValues: xAxis,
			YValues: yPrimary,
			//XValueFormatter: xAxisFormat,
			Style: plot.Style{StrokeColor: plot.ColorBlack},
		},
		plot.ContinuousSeries{
			Name:    ySecondaryName,
			YAxis:   plot.YAxisSecondary,
			XValues: xAxis,
			YValues: ySecondary,
			//XValueFormatter: xAxisFormat,
			Style: plot.Style{StrokeColor: plot.ColorAlternateGray, StrokeDashArray: plot.DashedLine},
		},
	}
	return graph
}
