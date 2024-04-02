package blas

import (
	"fmt"
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/ta-lib/plot"
	"github.com/wcharczuk/go-chart/v2"
)

// MatchWedge 模式匹配 楔形
func MatchWedge(waves Waves) *Wedge {
	if waves.PeakCount < 2 || waves.ValleyCount < 2 {
		return nil
	}

	m := Wedge{}
	m.Digits = waves.Digits
	// 两个高点
	m.TopLeft = waves.Peaks[waves.PeakCount-2]
	m.TopRight = waves.Peaks[waves.PeakCount-1]
	// 两个低点
	m.BottomLeft = waves.Valleys[waves.ValleyCount-2]
	m.BottomRight = waves.Valleys[waves.ValleyCount-1]
	return &m
}

// Wedge 楔形
type Wedge struct {
	OperationalSignal
	TopLeft     num.DataPoint // 顶 - 左
	TopRight    num.DataPoint // 顶 - 右
	BottomLeft  num.DataPoint // 底 - 左
	BottomRight num.DataPoint // 底 - 右
}

func (this *Wedge) ExportSeries(sample DataSample) []chart.Series {
	top := chart.ContinuousSeries{
		Name:    "头",
		XValues: []float64{float64(this.TopLeft.X), float64(this.TopRight.X)},
		YValues: []float64{this.TopLeft.Y, this.TopRight.Y},
		Style:   chart.Style{StrokeColor: plot.ColorRed, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorRed, DotWidth: plot.DotWidth},
	}
	bottom := chart.ContinuousSeries{
		Name:    "底",
		XValues: []float64{float64(this.BottomLeft.X), float64(this.BottomRight.X)},
		YValues: []float64{this.BottomLeft.Y, this.BottomRight.Y},
		Style:   chart.Style{StrokeColor: plot.ColorGreen, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorGreen, DotWidth: plot.DotWidth},
	}
	// 支撑线
	p1 := this.BottomLeft.ToPoint()
	p2 := this.BottomRight.ToPoint()
	supportLine := num.CalculateLineEquation(p1, p2)
	fmt.Printf("supportLine=%+v\n", supportLine)
	supportRows := sample.Len() - this.BottomRight.X
	supportX := make([]float64, supportRows)
	supportY := make([]float64, supportRows)
	downFound := false
	downX := float64(0)
	downY := float64(0)
	for i := 0; i < supportRows; i++ {
		pos := i + this.BottomRight.X
		supportX[i] = float64(pos)
		supportY[i] = supportLine.Y(supportX[i])
		supportY[i] = num.Decimal(supportY[i], this.Digits)
		if !downFound && i > 0 && pos > this.BottomRight.X && sample.Current(pos-1) > supportY[i-1] && sample.Current(pos) < supportY[i] {
			downFound = true
			downX = supportX[i]
			downY = sample.Current(pos)
		}
	}
	// 压力线
	p1 = this.TopLeft.ToPoint()
	p2 = this.TopRight.ToPoint()
	pressureLine := num.CalculateLineEquation(p1, p2)
	fmt.Printf("pressureLine=%+v\n", pressureLine)
	pressureRows := sample.Len() - this.TopRight.X
	pressureX := make([]float64, pressureRows)
	pressureY := make([]float64, pressureRows)
	upFound := false
	upX := float64(0)
	upY := float64(0)
	for i := 0; i < pressureRows; i++ {
		pos := i + this.TopRight.X
		pressureX[i] = float64(pos)
		pressureY[i] = pressureLine.Y(pressureX[i])
		pressureY[i] = num.Decimal(pressureY[i], this.Digits)
		if i > 0 {
			fmt.Println(this.TopRight.X, i, pos, "|", sample.Current(pos-1), pressureY[i-1], sample.Current(pos), pressureY[i])
		}
		if !upFound && i > 0 && pos > this.TopRight.X && sample.Current(pos-1) < pressureY[i-1] && sample.Current(pos) > pressureY[i] {
			upFound = true
			upX = pressureX[i]
			upY = sample.Current(pos)
			fmt.Println("\t", "FOUND")
		}
	}
	pressure := chart.ContinuousSeries{
		Name:    "压力线",
		XValues: pressureX,
		YValues: pressureY,
		Style:   chart.Style{StrokeColor: plot.ColorRed, StrokeDashArray: plot.DashedLine},
	}
	support := chart.ContinuousSeries{
		Name:    "支撑线",
		XValues: supportX,
		YValues: supportY,
		Style:   chart.Style{StrokeColor: plot.ColorGreen, StrokeDashArray: plot.DashedLine},
	}
	list := []chart.Series{top, bottom, pressure, support}
	if upFound {
		this.Signal = TradingBuy
		name := "突破"
		operation := "买入"
		color := plot.ColorRed
		tendency := chart.ContinuousSeries{
			Name:    name,
			XValues: []float64{upX},
			YValues: []float64{upY},
			Style:   chart.Style{StrokeColor: color, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorBlue, DotWidth: plot.DotWidth},
		}
		signal := chart.AnnotationSeries{
			Annotations: []chart.Value2{
				{
					Label:  operation + sample.Time(int(upX)),
					XValue: upX,
					YValue: upY,
					Style: chart.Style{
						StrokeColor: color,
						StrokeWidth: 0,
					},
				},
			},
		}
		list = append(list, tendency, signal)
	}
	if downFound {
		this.Signal = TradingSell
		name := "跌破"
		operation := "卖出"
		color := plot.ColorGreen
		tendency := chart.ContinuousSeries{
			Name:    name,
			XValues: []float64{downX},
			YValues: []float64{downY},
			Style:   chart.Style{StrokeColor: color, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorBlue, DotWidth: plot.DotWidth},
		}
		signal := chart.AnnotationSeries{
			Annotations: []chart.Value2{
				{
					Label:  operation + sample.Time(int(downX)),
					XValue: downX,
					YValue: downY,
					Style: chart.Style{
						StrokeColor: color,
						StrokeWidth: 0,
					},
				},
			},
		}
		list = append(list, tendency, signal)
	}
	return list
}
