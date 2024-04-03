package blas

import (
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/ta-lib/plot"
)

const (
	DoubleBottom = "W底"
	DoubleTop    = "M头"
)

// Triangle 三角形, W底和M头
//
// 双底的颈线有两个, 水平方向的TOP以及TOP于左右连线的平行线, 上升空间为top到底部区间的2倍
type Triangle struct {
	OperationalSignal
	Top         num.DataPoint // 顶或者底作为颈线
	Left        num.DataPoint // 左
	Right       num.DataPoint // 右
	PremiumRate float64       // 交易信号触发后的溢价率
	Name        string        // 形态名称
}

// IsUp 上W, 下M
func (this *Triangle) IsUp() bool {
	// TopValue > left
	a := this.Top.Y > this.Left.Y
	// TopValue > right
	b := this.Top.Y > this.Right.Y
	return a && b
}

// Space 计算溢价空间
func (this *Triangle) Space(periods int) {
	this.v2Space(periods)
}

func (this *Triangle) v1Space(periods int) {
	left := this.Left.ToPoint()
	right := this.Right.ToPoint()
	line := num.CalculateLineEquation(left, right)
	top := this.Top.ToPoint()
	//distance:= line.VerticalDistance(top)
	y := line.Y(top.X)
	distance := y - top.Y
	space := top.Y - distance
	this.PremiumRate = num.NetChangeRate(top.Y, space)
	_ = periods
}

func (this *Triangle) v2Space(periods int) {
	left := this.Left.ToPoint()
	right := this.Right.ToPoint()
	line := num.CalculateLineEquation(left, right)
	top := this.Top.ToPoint()
	// 1. 颈线
	neckLine := line.ParallelLine(top)
	// 2. 目标线
	targetLine := line.SymmetricParallelLine(top)
	neck := neckLine.Y(float64(periods))
	target := targetLine.Y(float64(periods))
	this.PremiumRate = num.NetChangeRate(neck, target)
}

// Fit 拟合
//
// 输出支撑线, 颈线, 压力线
func (this *Triangle) Fit() (neckLine, supportLine, pressureLine num.Line) {
	// 0. 确定左右连线和顶点
	left := this.Left.ToPoint()
	right := this.Right.ToPoint()
	line := num.CalculateLineEquation(left, right)
	top := this.Top.ToPoint()
	// 1. 颈线
	neckLine = line.ParallelLine(top)
	doubleLine := line.SymmetricParallelLine(top)
	// 2. 支撑线
	if this.IsUp() {
		// W底, 支撑线就是左右连线
		supportLine = line
	} else {
		// M头, 支撑线是line与top两倍距离的平行线
		supportLine = doubleLine
	}
	// 3. 压力线
	if this.IsUp() {
		// W底
		pressureLine = doubleLine
	} else {
		// M头
		pressureLine = line
	}
	// 4. 确定股价波动空间

	return
}

func (this *Triangle) Analyze(data []float64) []num.LinerTrend {
	neckLine, _, _ := this.Fit()
	offset := this.Right.X
	count := len(data)
	length := count - offset
	x := make([]float64, length)
	y := make([]float64, length)
	for i := 0; i < length; i++ {
		x[i] = float64(i + offset)
		y[i] = neckLine.Y(x[i])
		y[i] = num.Decimal(y[i], this.Digits)
	}

	tendency := num.Cross(data[offset:], y)
	for i := 0; i < len(tendency); i++ {
		tendency[i].X += offset
	}
	return tendency
}

// Detect 检测最近的趋势
func (this *Triangle) Detect(data []float64) (neck, support, pressure float64, neckTendency, supportTendency, pressureTendency int) {
	neckLine, supportLine, pressureLine := this.Fit()
	nx, ny, nt := neckLine.Extend(data, this.Digits)
	sx, sy, st := supportLine.Extend(data, this.Digits)
	px, py, pt := pressureLine.Extend(data, this.Digits)

	nl := len(ny)
	neck = ny[nl-1]
	neckTendency = nt
	sl := len(sy)
	support = sy[sl-1]
	supportTendency = st
	pl := len(py)
	pressure = py[pl-1]
	pressureTendency = pt
	_ = nx
	_ = sx
	_ = px

	return
}

// Export 导出趋势线图表数据
func (this *Triangle) Export(data []float64) []plot.Series {
	neckLine, supportLine, pressureLine := this.Fit()
	nx, ny, _ := neckLine.Extend(data, this.Digits)
	sx, sy, _ := supportLine.Extend(data, this.Digits)
	px, py, _ := pressureLine.Extend(data, this.Digits)
	top := plot.ContinuousSeries{
		Name:    "头",
		XValues: []float64{float64(this.Top.X)},
		YValues: []float64{this.Top.Y},
		Style:   plot.Style{StrokeColor: plot.ColorRed, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorRed, DotWidth: plot.DotWidth},
	}
	bottom := plot.ContinuousSeries{
		Name:    "底",
		XValues: []float64{float64(this.Left.X), float64(this.Right.X)},
		YValues: []float64{this.Left.Y, this.Right.Y},
		Style:   plot.Style{StrokeColor: plot.ColorGreen, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorGreen, DotWidth: plot.DotWidth},
	}
	list := []plot.Series{top, bottom}
	neck := plot.ContinuousSeries{
		Name:    "颈线",
		XValues: nx,
		YValues: ny,
		Style:   plot.Style{StrokeColor: plot.ColorBlue, StrokeDashArray: plot.DashedLine},
	}
	defaultFontSize := 9.0
	labelNeck := plot.LastValueAnnotationSeries(neck)
	labelNeck.Style.FontSize = defaultFontSize
	support := plot.ContinuousSeries{
		Name:    "支撑线",
		XValues: sx,
		YValues: sy,
		Style:   plot.Style{StrokeColor: plot.ColorGreen, StrokeDashArray: plot.DashedLine},
	}
	labelSupport := plot.LastValueAnnotationSeries(support)
	labelSupport.Style.FontSize = defaultFontSize
	pressure := plot.ContinuousSeries{
		Name:    "压力线",
		XValues: px,
		YValues: py,
		Style:   plot.Style{StrokeColor: plot.ColorRed, StrokeDashArray: plot.DashedLine},
	}
	labelPressure := plot.LastValueAnnotationSeries(pressure)
	labelPressure.Style.FontSize = defaultFontSize
	list = append(list, neck, support, pressure, labelNeck, labelSupport, labelPressure)
	return list
}

func (this *Triangle) ExportSeries(sample DataSample) []plot.Series {
	count := sample.Len()
	top := plot.ContinuousSeries{
		Name:    "头",
		XValues: []float64{float64(this.Top.X)},
		YValues: []float64{this.Top.Y},
		Style:   plot.Style{StrokeColor: plot.ColorRed, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorRed, DotWidth: plot.DotWidth},
	}
	bottom := plot.ContinuousSeries{
		Name:    "底",
		XValues: []float64{float64(this.Left.X), float64(this.Right.X)},
		YValues: []float64{this.Left.Y, this.Right.Y},
		Style:   plot.Style{StrokeColor: plot.ColorGreen, StrokeDashArray: plot.DashedLine, DotColor: plot.ColorGreen, DotWidth: plot.DotWidth},
	}
	p1 := this.Left.ToPoint()
	p2 := this.Right.ToPoint()
	bottomLine := num.CalculateLineEquation(p1, p2)
	topLine := bottomLine.ParallelLine(this.Top.ToPoint())
	neckX := make([]float64, count-this.Top.X)
	neckY := make([]float64, count-this.Top.X)
	for i := 0; i < count-this.Top.X; i++ {
		neckX[i] = float64(i + this.Top.X)
		neckY[i] = topLine.Y(neckX[i])
		neckY[i] = num.Decimal(neckY[i], this.Digits)
	}
	neck := plot.ContinuousSeries{
		Name:    "颈线",
		XValues: neckX,
		YValues: neckY,
		Style:   plot.Style{StrokeColor: plot.ColorBlue, StrokeDashArray: plot.DashedLine},
	}
	return []plot.Series{top, bottom, neck}
}

// MatchTriangle 模式匹配 三角形
func MatchTriangle(waves Waves) *Triangle {
	if waves.PeakCount+waves.ValleyCount < 3 {
		return nil
	}
	// 从后往前推导
	if waves.Peaks[waves.PeakCount-1].X < waves.Valleys[waves.ValleyCount-1].X {
		// 如果波峰的x值小于波谷的x值, 那么是底的模式
		// W底, 1个波峰和2个波谷
		if waves.ValleyCount >= 2 {
			var m Triangle
			m.Name = DoubleBottom
			m.Digits = waves.Digits
			m.Top = waves.Peaks[waves.PeakCount-1]
			m.Left = waves.Valleys[waves.ValleyCount-2]
			m.Right = waves.Valleys[waves.ValleyCount-1]
			m.Space(waves.Len())
			return &m
		}
	} else {
		// 如果波峰x值大于波谷的x值, 那么是顶的模式
		// M头, 2个波峰和1个波谷
		if waves.PeakCount >= 2 {
			var m Triangle
			m.Name = DoubleTop
			m.Digits = waves.Digits
			m.Top = waves.Valleys[waves.ValleyCount-1]
			m.Left = waves.Peaks[waves.PeakCount-2]
			m.Right = waves.Peaks[waves.PeakCount-1]
			m.Space(waves.Len())
			return &m
		}
	}
	return nil
}
