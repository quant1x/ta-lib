package algorithms

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pkg/chart"
	"gitee.com/quant1x/pkg/chart/drawing"
	"gitee.com/quant1x/pkg/testify/assert"
	"gitee.com/quant1x/ta-lib/plot"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"
)

const (
	lineChartXAxisName  = "Date"
	lineChartYAxisName  = "Count"
	lineChartHeight     = 700
	lineChartWidth      = 1280
	colorMultiplier     = 256
	imgStrPrefix        = "data:image/png;base64,"
	pieLabelFormat      = "%v %v"
	barChartTryAgainErr = "invalid data range; cannot be zero"
)

var (
	lineChartStyle = chart.Style{
		Padding: chart.Box{
			Top:  30,
			Left: 150,
		},
	}

	defaultChartStyle = chart.Style{
		Padding: chart.Box{
			Top: 30,
		},
	}

	timeFormat = chart.TimeDateValueFormatter
)

type LineYValue struct {
	Name   string
	Values []float64
}

type ChartValue struct {
	Name  string
	Value float64
}

// createLineChart 创建线性图
func createLineChart(title string, endTime time.Time, values []LineYValue) (img string, err error) {
	if len(values) == 0 {
		return
	}
	// 1、计算X轴
	lenX := len(values[0].Values)
	// X轴内容xValues 及 X轴坐标ticks
	var xValues []time.Time
	var ticks []chart.Tick
	for i := lenX - 1; i >= 0; i-- {
		curTime := endTime.AddDate(0, 0, -i)
		xValues = append(xValues, curTime)
		ticks = append(ticks, chart.Tick{Value: getNsec(curTime), Label: timeFormat(curTime)})
	}

	// 2、生成Series
	var series []chart.Series
	for _, yValue := range values {
		series = append(series, chart.TimeSeries{
			Name: yValue.Name,
			Style: chart.Style{
				// 随机渲染线条颜色
				StrokeColor: drawing.Color{
					R: uint8(rand.Intn(colorMultiplier)),
					G: uint8(rand.Intn(colorMultiplier)),
					B: uint8(rand.Intn(colorMultiplier)),
					A: uint8(colorMultiplier - 1), // 透明度
				},
			},
			XValues: xValues,
			YValues: yValue.Values,
		})
	}

	// 3、新建图形
	graph := chart.Chart{
		Title:      title,
		Background: lineChartStyle,
		Width:      lineChartWidth,
		Height:     lineChartHeight,
		XAxis: chart.XAxis{
			Name:           lineChartXAxisName,
			ValueFormatter: timeFormat,
			Ticks:          ticks,
		},
		YAxis: chart.YAxis{
			Name: lineChartYAxisName,
		},
		Series: series,
	}
	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}

	// 4、输出目标
	img, err = writeLineChart(&graph)

	return
}

// getNsec 获取纳秒数
func getNsec(cur time.Time) float64 {
	return float64(cur.Unix() * int64(time.Second))
}

func writeLineChartToPng(c *chart.Chart) (img string, err error) {
	f, _ := os.Create("graph.png")
	err = c.Render(chart.PNG, f)
	return
}

func writeLineChart(c *chart.Chart) (img string, err error) {
	var imgContent bytes.Buffer
	err = c.Render(chart.PNG, &imgContent)
	if err != nil {
		return
	}

	os.WriteFile("1.png", imgContent.Bytes(), 0644)

	img = imgToStr(imgContent)
	return
}

func imgToStr(imgContent bytes.Buffer) string {
	return imgStrPrefix + base64.StdEncoding.EncodeToString(imgContent.Bytes())
}

// createPieChart 创建饼图
func createPieChart(title string, pieValues []ChartValue) (img string, err error) {
	if len(pieValues) == 0 {
		return
	}
	// 1、构建value
	var values []chart.Value
	for _, v := range pieValues {

		values = append(values, chart.Value{
			Value: v.Value,
			Label: fmt.Sprintf(pieLabelFormat, getSimpleSensType(v.Name), formatValue(v.Value)),
		})
	}

	// 2、新建饼图
	pie := chart.PieChart{
		Title:      title,
		Background: defaultChartStyle,
		Values:     values,
	}

	// 4、输出目标
	img, err = writePieChart(&pie)

	return
}

func formatValue(f float64) string {
	return fmt.Sprintf("%.2fW", f/10000)
}

func getSimpleSensType(name string) string {
	if name == "个人数据" {
		return "Personal"
	}
	return "Other"
}

func writePieChartToPng(c *chart.PieChart) (img string, err error) {
	f, _ := os.Create("pie.png")
	err = c.Render(chart.PNG, f)
	return
}

func writePieChart(c *chart.PieChart) (img string, err error) {
	var imgContent bytes.Buffer
	err = c.Render(chart.PNG, &imgContent)
	if err != nil {
		return
	}

	img = imgToStr(imgContent)
	return
}

// createBarChart 创建柱状图
func createBarChart(title string, barValues []ChartValue) (img string, err error) {
	if len(barValues) == 0 {
		return
	}
	// 1、构建value
	var values []chart.Value
	for _, v := range barValues {
		values = append(values, chart.Value{
			Value: v.Value,
			Label: v.Name,
		})
	}

	// 2、新建饼图
	bar := chart.BarChart{
		XAxis: chart.Style{
			TextWrap: 0, // default 1为可以溢出规定的范围
		},
		Width:      2560,
		BarWidth:   50,
		BarSpacing: 300,
		Title:      title,
		Background: defaultChartStyle,
		Bars:       values,
	}

	// 4、输出目标
	img, err = writeBarChart(&bar)
	if err != nil && err.Error() == barChartTryAgainErr {
		// 添加一个隐藏条目，设置透明度A为0, 设置任意属性如R不为0即可
		values = append(values, chart.Value{
			Style: chart.Style{
				StrokeColor: drawing.Color{R: 1},
			},
			Value: 0,
			Label: "",
		})
		bar.Bars = values
		img, err = writeBarChart(&bar)
	}

	return
}

func writeBarChartToPng(c *chart.BarChart) (img string, err error) {
	f, _ := os.Create("bar.png")
	err = c.Render(chart.PNG, f)
	return
}

func writeBarChart(c *chart.BarChart) (img string, err error) {
	var imgContent bytes.Buffer
	err = c.Render(chart.PNG, &imgContent)
	if err != nil {
		return
	}

	img = imgToStr(imgContent)
	return
}

func TestCreateLineChart(t *testing.T) {
	testAssert := assert.New(t)

	tests := []struct {
		title     string
		endTime   time.Time
		barValues []LineYValue
	}{
		{"line chart", time.Now(), []LineYValue{
			{"asd", []float64{1, 2, 300, 100, 200, 6, 700}},
			{"hgj", []float64{400, 500000, 200, 50, 5, 800, 7}},
			{"dfg45r", []float64{1, 2, 700, 100, 200, 6, 700}},
			{"2342sr", []float64{400, 500000, 200, 50, 5, 800, 7}},
			{"das21-asd", []float64{300000, 200000, 400000, 100000, 400000, 450000, 400000}},
			{"csc", []float64{400, 500000, 200, 50, 5, 800, 7}},
			{"mhj", []float64{1, 2, 300, 100, 200, 6, 700}},
			{"876ijgh", []float64{400, 500000, 200, 50, 5, 800, 7}},
			{"fbfdv", []float64{1, 2, 300, 100, 200, 6, 700}},
			{"67ds", []float64{400, 10000, 200, 50, 5, 800, 7}},
			{"67bdfv", []float64{1, 2, 300, 100, 200, 6, 700}},
			{"sdf324", []float64{400, 500000, 200, 50, 5, 800, 7}},
			{"vdf67", []float64{1, 2, 300, 100, 200, 6, 700}},
			{"vdfs234", []float64{400, 500000, 200, 50, 5, 800, 7}},
			{"123sdf", []float64{1, 2, 700, 100, 200, 6, 700}},
			{"aasdasd", []float64{400, 500000, 200, 50, 5, 800, 7}},
			{"aasd", []float64{1, 2, 300, 100, 200, 6, 700}},
			{"basd", []float64{400, 500000, 200, 50, 5, 800, 7}},
			{"cczx", []float64{1, 2, 300, 100, 200, 6, 700}},
			{"qweqw", []float64{400, 500000, 200, 50, 5, 800, 7}},
			{"asdadf", []float64{1, 2, 300, 100, 200, 6, 700}},
			{"fghfh", []float64{400, 500000, 200, 50, 5, 800, 7}},
			{"erttyrt", []float64{1, 2, 300, 100, 200, 6, 700}}}},
	}

	for _, test := range tests {
		img, err := createLineChart(test.title, test.endTime, test.barValues)
		fmt.Println(img)
		testAssert.Equal(img, "")
		testAssert.Equal(err, nil)
	}
}
func TestFindPeekV1(t *testing.T) {
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
	data = []float64{13648.136043,
		2439.662155,
		9298.772911,
		17962.721514,
		20266.670986,
		29784.827490,
		31288.613046,
		2995.089472,
		21175.518221,
		19288.732784,
		15162.171623,
		11400.202960,
		15218.510132,
		21472.704585,
		30708.646304,
		2850.696496,
		25527.758715,
		18904.823877,
		29190.296136,
		3236.138204,
		9385.428575,
		1650.822572,
		16132.101779,
		24588.024061,
		9055.532039,
		7211.527865,
		10601.268564,
		14504.043245,
		23596.308422,
		16316.340951,
		12240.454139,
		28975.462939,
		30888.247967,
		20104.555282,
		4224.997442,
		19785.656907,
		12256.643330,
		31718.199922,
		9178.473112,
		22598.125747,
		17126.588720,
		6182.753235,
		22061.923230,
		1614.560101,
		5901.891702,
		1461.066351,
		28321.328065,
		26873.420346,
		3780.967845,
		13133.275318,
		3847.314880,
		18306.960712,
		30380.484308,
		8204.325114,
		31675.693925,
		11193.857313,
		6672.715965,
		21306.469121,
		31147.030615,
		19926.394289,
		2033.212562,
		11952.305064,
		5320.049134,
		7400.898499,
		1670.676265,
		28856.202135,
		25385.334397,
		11936.458395,
		26625.750572,
		24122.705811,
		19899.619112,
		12610.981787,
		11496.894766,
		2843.278749,
		10933.664116,
		17557.459010,
		14737.518755,
		20654.475950,
		16432.666921,
		26061.646994,
		3109.862592,
		14838.833941,
		18874.159667,
		5989.506257,
		19562.564692,
		1662.143404,
		18423.262351,
		26955.030459,
		15991.231156,
		14048.796405,
		4769.826664,
		904.942429,
		24213.433215,
		25475.399550,
		9393.779021,
		3686.615873,
		12002.932192,
		26524.599842,
		26936.850623,
		21287.630450,
		30724.479228,
		30179.772551,
		3606.383756,
		20745.198375,
		15385.729654,
		2128.662142,
		28728.679217,
		15911.364126,
		24681.709723,
		1931.591342,
		8398.638877,
		20834.219798,
		4275.325134,
		20433.463578,
		12318.179039,
		24502.340022,
		20893.301829,
		12207.650089,
		9600.594648,
		10884.465776,
		29405.649175,
		14600.538026,
		14159.902483,
		14533.946228,
		30249.024294,
		7011.806482,
		28236.889470,
		636.013530,
		10936.476264,
		24512.878262,
		10969.718449,
		19801.804933,
		14496.677187,
		325.203284,
		19170.597359,
		19250.187147,
		20781.358431,
		10967.060243,
		15785.576178,
		22456.766248,
		28409.683246,
		1761.851077,
		3147.582417,
		20793.049135,
		24450.268520,
		31614.686131}
	pv := InitPV(data)
	pv.Find()
	fmt.Println(pv)
	// 输出图表
	rows = len(data)
	dataY := data
	x := num.Range[float64](rows)
	dataX := x
	font, _ := plot.GetDefaultFont()
	peekY := make([]float64, 0, pv.Pcnt)
	peekX := make([]float64, 0, pv.Pcnt)
	frontValue := float64(0)
	for i, v := range pv.PosPeak {
		if i >= pv.Pcnt {
			break
		}
		peekX = append(peekX, float64(v))
		peekY = append(peekY, pv.Data[v])
	}
	valleyX := make([]float64, 0, pv.Vcnt)
	valleyY := make([]float64, 0, pv.Vcnt)
	for i, v := range pv.PosValley {
		if i >= pv.Pcnt {
			break
		}
		valleyX = append(valleyX, float64(v))
		valleyY = append(valleyY, pv.Data[v])
	}
	_ = frontValue
	graph := chart.Chart{
		Title: "走势图",
		Font:  font,
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: dataX,
				YValues: dataY,
				//XValueFormatter: func(v interface{}) string {
				//	f := v.(float64)
				//	idx := int(f)
				//	return DATE[idx]
				//},
				Style: chart.Style{StrokeColor: chart.ColorBlack},
			},
			chart.ContinuousSeries{
				Name:    "波峰",
				XValues: peekX,
				YValues: peekY,
				//XValueFormatter: func(v interface{}) string {
				//	f := v.(float64)
				//	idx := int(f)
				//	return DATE[idx]
				//},
				Style: chart.Style{StrokeColor: chart.ColorRed},
			},
			chart.ContinuousSeries{
				Name:    "波谷",
				XValues: valleyX,
				YValues: valleyY,
				//XValueFormatter: func(v interface{}) string {
				//	f := v.(float64)
				//	idx := int(f)
				//	return DATE[idx]
				//},
				Style: chart.Style{StrokeColor: chart.ColorGreen},
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
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
	_ = DATE
}

func TestFindPeek(t *testing.T) {
	code := "sh000001"
	//code = "sh603933"
	//code = "sh603230"
	//code = "sh605577"
	date := "20240103"
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
	pv := InitPV(data)
	pv.Find()
	fmt.Println(pv)
	// 输出图表
	rows = len(data)
	dataY := data
	x := num.Range[float64](rows)
	dataX := x
	font, _ := plot.GetDefaultFont()
	font.Bounds(1)
	peekY := make([]float64, 0, pv.Pcnt)
	peekX := make([]float64, 0, pv.Pcnt)
	frontValue := float64(0)
	fmt.Println("波峰")
	for i, v := range pv.PosPeak {
		if i >= pv.Pcnt {
			break
		}
		fmt.Println("\t", i, DATE[v], pv.Data[v])
		peekX = append(peekX, float64(v))
		peekY = append(peekY, pv.Data[v])
	}
	fmt.Println("波谷")
	valleyX := make([]float64, 0, pv.Vcnt)
	valleyY := make([]float64, 0, pv.Vcnt)
	for i, v := range pv.PosValley {
		if i >= pv.Vcnt {
			break
		}
		fmt.Println("\t", i, DATE[v], pv.Data[v])
		valleyX = append(valleyX, float64(v))
		valleyY = append(valleyY, pv.Data[v])
	}
	// 找到最后2个波谷
	vn := len(valleyX)
	leftX := valleyX[vn-2]
	leftY := valleyY[vn-2]
	rightX := valleyX[vn-1]
	rightY := valleyY[vn-1]
	fmt.Println("最近的两个底部")
	fmt.Println("\t=> 左:", leftX, leftY)
	fmt.Println("\t=> 右:", rightX, rightY)
	fmt.Println("计算斜率")
	xl := num.Slope(int(leftX), leftY, int(rightX), rightY)
	fmt.Println("\t=>", xl)
	// 找到最后一个波峰
	fmt.Println("最近的波峰")
	pn := len(peekX)
	// 颈线的点
	neckPointX := int(peekX[pn-1])
	neckPointY := peekY[pn-1]
	fmt.Println("\t颈线:", neckPointX, neckPointY)
	fmt.Println("计算目前为止")
	//cjx := num.TriangleBevel(xl, neckPointX, neckPointY, rows-neckPointX)
	supportLine := calculateLineEquation(Point{x: leftX, y: leftY}, Point{x: rightX, y: rightY})
	fmt.Println("supportLine =", supportLine)
	neckLine := calculateEquidistantLine(supportLine, Point{x: float64(neckPointX), y: neckPointY})
	fmt.Println("neckLine =", neckLine)
	fmt.Println("\t目前颈线所在位置", neckLine.m*float64(rows-1)+neckLine.c)
	pressurePointX := neckPointX
	supportY := supportLine.m*float64(pressurePointX) + supportLine.c
	d := neckPointY - supportY
	pressurePointY := neckPointY + d
	pressureLine := calculateEquidistantLine(neckLine, Point{x: float64(pressurePointX), y: pressurePointY})
	fmt.Println("pressureLine =", pressureLine)
	//high := pressureLine.m*float64(rows-1) + pressureLine.c
	neckX := []float64{float64(neckPointX), float64(rows - 1)}
	neckY := []float64{neckLine.m*float64(neckPointX) + neckLine.c, neckLine.m*float64(rows-1) + neckLine.c}
	pressureX := []float64{float64(neckPointX), float64(rows - 1)}
	pressureY := []float64{pressurePointY, pressureLine.m*float64(rows-1) + pressureLine.c}
	fmt.Println("\t最近压力", pressureY)
	fmt.Println("\t明天压力", pressureLine.m*float64(rows)+pressureLine.c)
	neckValue := neckPointY - math.Abs(leftY-rightY)
	pressureMin := neckValue*2 - rightY
	pressureMax := neckValue*2 - leftY
	if pressureMin > pressureMax {
		pressureMin, pressureMax = pressureMax, pressureMin
	}
	fmt.Println("\t反弹高度", pressureMin, "~", pressureMax)
	_ = frontValue

	xAxisFormat := func(v interface{}) string {
		f := v.(float64)
		idx := int(f)
		return DATE[idx]
	}
	lineChartStyle := chart.Style{
		Padding: chart.Box{
			//Top:  10,
			Left: 50,
		},
	}
	jx := chart.ContinuousSeries{
		Name:            "颈线",
		XValues:         neckX,
		YValues:         neckY,
		XValueFormatter: xAxisFormat,
		Style:           chart.Style{StrokeColor: chart.ColorBlue},
	}
	ylx := chart.ContinuousSeries{
		Name:            "压力线",
		XValues:         pressureX,
		YValues:         pressureY,
		XValueFormatter: xAxisFormat,
		Style:           chart.Style{StrokeColor: chart.ColorRed, StrokeDashArray: []float64{5.0, 5.0}},
	}
	graph := chart.Chart{
		Title:      code + "走势图",
		Font:       font,
		Background: lineChartStyle,
		//Width:      lineChartWidth,
		//Height:     lineChartHeight,
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:            "CLOSE",
				XValues:         dataX,
				YValues:         dataY,
				XValueFormatter: xAxisFormat,
				Style:           chart.Style{StrokeColor: chart.ColorBlack},
			},
			chart.ContinuousSeries{
				Name:            "波峰",
				XValues:         peekX,
				YValues:         peekY,
				XValueFormatter: xAxisFormat,
				Style:           chart.Style{StrokeColor: chart.ColorRed, DotWidth: 5, DotColor: chart.ColorRed},
			},
			chart.ContinuousSeries{
				Name:            "波谷",
				XValues:         valleyX,
				YValues:         valleyY,
				XValueFormatter: xAxisFormat,
				Style:           chart.Style{StrokeColor: chart.ColorGreen, DotWidth: 5, DotColor: chart.ColorGreen},
			},
			jx,
			chart.LastValueAnnotationSeries(jx),
			ylx,
			chart.LastValueAnnotationSeries(ylx),
		},
	}

	//_ = pressureX
	//_ = pressureY

	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
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
	//_ = plot.OpenImage(graph)
	_ = DATE
}
