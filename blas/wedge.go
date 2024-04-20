package blas

import (
	"fmt"
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/ta-lib/plot"
	"slices"
)

// MatchWedge 模式匹配 楔形
func MatchWedge(waves Waves) *Wedge {
	return v1MatchWedge(waves)
}

func v1MatchWedge(waves Waves) *Wedge {
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

func v2MatchWedge(waves Waves) *Wedge {
	if waves.PeakCount < 2 || waves.ValleyCount < 2 {
		return nil
	}

	m := Wedge{}
	m.Digits = waves.Digits
	peaks := slices.Clone(waves.Peaks)
	slices.SortFunc(peaks, func(a, b num.DataPoint) int {
		return Desc(a.Y, b.Y)
	})
	peak1 := peaks[0]
	peak2 := peaks[1]
	if peak1.X < peak2.X {
		peak1, peak2 = peak2, peak1
	}
	valleys := slices.Clone(waves.Valleys)
	slices.SortFunc(valleys, func(a, b num.DataPoint) int {
		return Asc(a.Y, b.Y)
	})
	valley1 := valleys[0]
	valley2 := valleys[1]
	if valley1.X < valley2.X {
		valley1, valley2 = valley2, valley1
	}

	// 两个高点
	m.TopLeft = peak1
	m.TopRight = peak2
	// 两个低点
	m.BottomLeft = valley1
	m.BottomRight = valley2
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

func (this *Wedge) Cross() (dp, top, bottom num.DataPoint, ok bool) {
	line1 := num.CalculateLineEquation(this.TopLeft.ToPoint(), this.TopRight.ToPoint())
	line2 := num.CalculateLineEquation(this.BottomLeft.ToPoint(), this.BottomRight.ToPoint())
	m1, _, b1 := line1.Equation()
	m2, _, b2 := line2.Equation()
	if m1 == m2 {
		ok = false
		return
	}
	// 计算交点的x坐标
	x := (b2 - b1) / (m1 - m2)
	// 计算交点的y坐标
	y := m1*x + b1
	n := int(x)
	dp.X = n
	if x > float64(n) {
		dp.X++
	}
	dp.Y = y
	ok = true
	top.X = dp.X
	top.Y = line1.Y(float64(top.X))
	bottom.X = dp.X
	bottom.Y = line2.Y(float64(bottom.X))
	return
}

// 双顶
func (this *Wedge) doubleTop() *num.DataPoint {
	if this.TopLeft.X < this.BottomLeft.X && this.BottomLeft.X < this.TopRight.X {
		return &this.BottomLeft
	}

	if this.TopLeft.X < this.BottomRight.X && this.BottomRight.X < this.TopRight.X {
		return &this.BottomRight
	}
	return nil
}

// 双底
func (this *Wedge) doubleBottom() *num.DataPoint {
	if this.BottomLeft.X < this.TopLeft.X && this.TopLeft.X < this.BottomRight.X {
		return &this.TopLeft
	}

	if this.BottomLeft.X < this.TopRight.X && this.TopRight.X < this.BottomRight.X {
		return &this.TopRight
	}
	return nil
}

// 锚点
func (this *Wedge) anchorPoints() []num.DataPoint {
	var list []num.DataPoint
	anchorPoint := this.doubleTop()
	if anchorPoint != nil {
		list = append(list, *anchorPoint)
	}
	anchorPoint = this.doubleBottom()
	if anchorPoint != nil {
		list = append(list, *anchorPoint)
	}
	return list
}

func (this *Wedge) Fit() (neckLine, supportLine, pressureLine num.Line) {
	pressureLine = num.CalculateLineEquation(this.TopLeft.ToPoint(), this.TopRight.ToPoint())
	supportLine = num.CalculateLineEquation(this.BottomLeft.ToPoint(), this.BottomRight.ToPoint())
	anchorPoint := this.doubleTop()
	if anchorPoint == nil {
		anchorPoint = this.doubleBottom()
	}
	if anchorPoint != nil {
		neckLine = pressureLine.ParallelLine(anchorPoint.ToPoint())
	}
	return
}

func (this *Wedge) NeckLines() []num.Line {
	var list []num.Line
	pressureLine := num.CalculateLineEquation(this.TopLeft.ToPoint(), this.TopRight.ToPoint())
	supportLine := num.CalculateLineEquation(this.BottomLeft.ToPoint(), this.BottomRight.ToPoint())
	anchorPoint := this.doubleTop()
	if anchorPoint != nil {
		neckLine := pressureLine.ParallelLine(anchorPoint.ToPoint())
		list = append(list, neckLine)
	}
	anchorPoint = this.doubleBottom()
	if anchorPoint != nil {
		neckLine := supportLine.ParallelLine(anchorPoint.ToPoint())
		list = append(list, neckLine)
	}
	return list
}

func (this *Wedge) getChartSeries(sample DataSample, line num.Line, point num.DataPoint) plot.Series {
	count := sample.Len()
	neckX := make([]float64, count-point.X)
	neckY := make([]float64, count-point.X)
	for i := 0; i < count-point.X; i++ {
		neckX[i] = float64(i + point.X)
		neckY[i] = line.Y(neckX[i])
		neckY[i] = num.Decimal(neckY[i], this.Digits)
	}
	neck := plot.ContinuousSeries{
		Name:    "颈线",
		XValues: neckX,
		YValues: neckY,
		Style:   plot.Style{StrokeColor: plot.ColorBlue, StrokeDashArray: plot.DashedLine},
	}
	return neck
}

func (this *Wedge) NeckSeries(sample DataSample) []plot.Series {
	var list []plot.Series
	pressureLine := num.CalculateLineEquation(this.TopLeft.ToPoint(), this.TopRight.ToPoint())
	supportLine := num.CalculateLineEquation(this.BottomLeft.ToPoint(), this.BottomRight.ToPoint())
	anchorPoint := this.doubleTop()
	if anchorPoint != nil {
		neckLine := pressureLine.ParallelLine(anchorPoint.ToPoint())
		neckSeries := this.getChartSeries(sample, neckLine, *anchorPoint)
		list = append(list, neckSeries)
	}
	anchorPoint = this.doubleBottom()
	if anchorPoint != nil {
		neckLine := supportLine.ParallelLine(anchorPoint.ToPoint())
		neckSeries := this.getChartSeries(sample, neckLine, *anchorPoint)
		list = append(list, neckSeries)
	}
	return list
}

func (this *Wedge) ExportSeries(sample DataSample) []plot.Series {
	return this.v2ExportSeries(sample)
}

func (this *Wedge) v2ExportSeries(sample DataSample) []plot.Series {
	top := plot.ContinuousSeries{
		Name:    "头",
		XValues: []float64{float64(this.TopLeft.X), float64(this.TopRight.X)},
		YValues: []float64{this.TopLeft.Y, this.TopRight.Y},
		Style:   plot.Style{StrokeColor: plot.ColorRed, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorRed, DotWidth: plot.DotWidth},
	}
	bottom := plot.ContinuousSeries{
		Name:    "底",
		XValues: []float64{float64(this.BottomLeft.X), float64(this.BottomRight.X)},
		YValues: []float64{this.BottomLeft.Y, this.BottomRight.Y},
		Style:   plot.Style{StrokeColor: plot.ColorGreen, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorGreen, DotWidth: plot.DotWidth},
	}
	// 支撑线
	p1 := this.BottomLeft.ToPoint()
	p2 := this.BottomRight.ToPoint()
	supportLine := num.CalculateLineEquation(p1, p2)
	fmt.Printf("supportLine=%+v\n", supportLine)
	supportRows := sample.Len() - this.BottomRight.X
	supportX := make([]float64, supportRows)
	supportY := make([]float64, supportRows)
	level2High := this.TopLeft.Y < this.TopRight.Y
	level2Low := this.BottomLeft.Y > this.BottomRight.Y
	supportCondition := level2High && level2Low
	chaodiFound := false
	chaodiX := float64(0)
	chaodiY := float64(0)

	for i := 0; i < supportRows && supportCondition; i++ {
		pos := i + this.BottomRight.X
		if !chaodiFound && i > 0 && pos > this.BottomRight.X && sample.Current(pos) > sample.Current(pos-1) {
			chaodiFound = true
			chaodiX = float64(pos)
			chaodiY = sample.Current(pos)
			fmt.Println("\t抄底:", "FOUND")
		}
	}
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
			fmt.Println("\tS:", "FOUND")
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
			fmt.Println("\tB:", "FOUND")
		}
	}
	pressure := plot.ContinuousSeries{
		Name:    "压力线",
		XValues: pressureX,
		YValues: pressureY,
		Style:   plot.Style{StrokeColor: plot.ColorRed, StrokeDashArray: plot.DashedLine},
	}
	support := plot.ContinuousSeries{
		Name:    "支撑线",
		XValues: supportX,
		YValues: supportY,
		Style:   plot.Style{StrokeColor: plot.ColorGreen, StrokeDashArray: plot.DashedLine},
	}
	list := []plot.Series{top, bottom, pressure, support}
	if chaodiFound {
		this.Signal = TradingBuy
		name := "2低"
		operation := "2低"
		color := plot.ColorRed
		tendency := plot.ContinuousSeries{
			Name:    name,
			XValues: []float64{chaodiX},
			YValues: []float64{chaodiY},
			Style:   plot.Style{StrokeColor: color, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorBlue, DotWidth: plot.DotWidth},
		}
		signal := plot.AnnotationSeries{
			Annotations: []plot.Value2{
				{
					Label:  operation + sample.Time(int(chaodiX)),
					XValue: chaodiX,
					YValue: chaodiY,
					Style: plot.Style{
						StrokeColor: color,
						StrokeWidth: 0,
					},
				},
			},
		}
		list = append(list, tendency, signal)
	}
	if upFound {
		this.Signal = TradingBuy
		name := "突破"
		operation := "买入"
		color := plot.ColorRed
		tendency := plot.ContinuousSeries{
			Name:    name,
			XValues: []float64{upX},
			YValues: []float64{upY},
			Style:   plot.Style{StrokeColor: color, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorBlue, DotWidth: plot.DotWidth},
		}
		signal := plot.AnnotationSeries{
			Annotations: []plot.Value2{
				{
					Label:  operation + sample.Time(int(upX)),
					XValue: upX,
					YValue: upY,
					Style: plot.Style{
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
		tendency := plot.ContinuousSeries{
			Name:    name,
			XValues: []float64{downX},
			YValues: []float64{downY},
			Style:   plot.Style{StrokeColor: color, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorBlue, DotWidth: plot.DotWidth},
		}
		signal := plot.AnnotationSeries{
			Annotations: []plot.Value2{
				{
					Label:  operation + sample.Time(int(downX)),
					XValue: downX,
					YValue: downY,
					Style: plot.Style{
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

func (this *Wedge) v1ExportSeries(sample DataSample) []plot.Series {
	top := plot.ContinuousSeries{
		Name:    "头",
		XValues: []float64{float64(this.TopLeft.X), float64(this.TopRight.X)},
		YValues: []float64{this.TopLeft.Y, this.TopRight.Y},
		Style:   plot.Style{StrokeColor: plot.ColorRed, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorRed, DotWidth: plot.DotWidth},
	}
	bottom := plot.ContinuousSeries{
		Name:    "底",
		XValues: []float64{float64(this.BottomLeft.X), float64(this.BottomRight.X)},
		YValues: []float64{this.BottomLeft.Y, this.BottomRight.Y},
		Style:   plot.Style{StrokeColor: plot.ColorGreen, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorGreen, DotWidth: plot.DotWidth},
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
			fmt.Println("\tS:", "FOUND")
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
			fmt.Println("\tB:", "FOUND")
		}
	}
	pressure := plot.ContinuousSeries{
		Name:    "压力线",
		XValues: pressureX,
		YValues: pressureY,
		Style:   plot.Style{StrokeColor: plot.ColorRed, StrokeDashArray: plot.DashedLine},
	}
	support := plot.ContinuousSeries{
		Name:    "支撑线",
		XValues: supportX,
		YValues: supportY,
		Style:   plot.Style{StrokeColor: plot.ColorGreen, StrokeDashArray: plot.DashedLine},
	}
	list := []plot.Series{top, bottom, pressure, support}
	if upFound {
		this.Signal = TradingBuy
		name := "突破"
		operation := "买入"
		color := plot.ColorRed
		tendency := plot.ContinuousSeries{
			Name:    name,
			XValues: []float64{upX},
			YValues: []float64{upY},
			Style:   plot.Style{StrokeColor: color, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorBlue, DotWidth: plot.DotWidth},
		}
		signal := plot.AnnotationSeries{
			Annotations: []plot.Value2{
				{
					Label:  operation + sample.Time(int(upX)),
					XValue: upX,
					YValue: upY,
					Style: plot.Style{
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
		tendency := plot.ContinuousSeries{
			Name:    name,
			XValues: []float64{downX},
			YValues: []float64{downY},
			Style:   plot.Style{StrokeColor: color, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorBlue, DotWidth: plot.DotWidth},
		}
		signal := plot.AnnotationSeries{
			Annotations: []plot.Value2{
				{
					Label:  operation + sample.Time(int(downX)),
					XValue: downX,
					YValue: downY,
					Style: plot.Style{
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
