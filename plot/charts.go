package plot

import (
	"bytes"
	"fmt"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/pkg/chart"
	"gitee.com/quant1x/pkg/chart/drawing"
	"os"
	_ "unsafe" // for go:linkname
)

const (
	DotWidth = 3  // 点的宽度
	XTickMax = 50 // X轴刻度多个不会重叠
)
const (
	YAxisPrimary   = chart.YAxisPrimary   // 主要坐标
	YAxisSecondary = chart.YAxisSecondary // 次要坐标
)

const (
	// TickPositionUnset 表示使用默认的刻度位置
	TickPositionUnset = chart.TickPositionUnset
	// TickPositionBetweenTicks 上一个刻度和当前刻度之间的刻度绘制标签.
	TickPositionBetweenTicks = chart.TickPositionBetweenTicks
	// TickPositionUnderTick 在记号下方绘制记号.
	TickPositionUnderTick = chart.TickPositionUnderTick
)

var (
	DashedLine = []float64{5.0, 5.0} // 虚线
)

var (
	ColorRed   = drawing.Color{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF} // 红色
	ColorGreen = drawing.Color{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF} // 绿色
	ColorBlue  = drawing.Color{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF} // 蓝色
	ColorBlack = drawing.Color{R: 51, G: 51, B: 51, A: 255}        // 黑色

	// ColorWhite is white.
	ColorWhite = drawing.Color{R: 255, G: 255, B: 255, A: 255}
	// ColorBlue is the basic theme blue color.
	//ColorBlue = drawing.Color{R: 0, G: 116, B: 217, A: 255}
	// ColorCyan is the basic theme cyan color.
	ColorCyan = drawing.Color{R: 0, G: 217, B: 210, A: 255}
	// ColorGreen is the basic theme green color.
	//ColorGreen = drawing.Color{R: 0, G: 217, B: 101, A: 255}
	// ColorRed is the basic theme red color.
	//ColorRed = drawing.Color{R: 217, G: 0, B: 116, A: 255}
	// ColorOrange is the basic theme orange color.
	ColorOrange = drawing.Color{R: 217, G: 101, B: 0, A: 255}
	// ColorYellow is the basic theme yellow color.
	ColorYellow = drawing.Color{R: 217, G: 210, B: 0, A: 255}
	// ColorBlack is the basic theme black color.
	//ColorBlack = drawing.Color{R: 51, G: 51, B: 51, A: 255}
	// ColorLightGray is the basic theme light gray color.
	ColorLightGray = drawing.Color{R: 239, G: 239, B: 239, A: 255}

	// ColorAlternateBlue is a alternate theme color.
	ColorAlternateBlue = drawing.Color{R: 106, G: 195, B: 203, A: 255}
	// ColorAlternateGreen is a alternate theme color.
	ColorAlternateGreen = drawing.Color{R: 42, G: 190, B: 137, A: 255}
	// ColorAlternateGray is a alternate theme color.
	ColorAlternateGray = drawing.Color{R: 110, G: 128, B: 139, A: 255}
	// ColorAlternateYellow is a alternate theme color.
	ColorAlternateYellow = drawing.Color{R: 240, G: 174, B: 90, A: 255}
	// ColorAlternateLightGray is a alternate theme color.
	ColorAlternateLightGray = drawing.Color{R: 187, G: 190, B: 191, A: 255}
	// ColorTransparent is a transparent (alpha zero) color.
	ColorTransparent = drawing.Color{R: 1, G: 1, B: 1, A: 0}
)

// 映射chart工具库, 收敛功能

type Series = chart.Series
type ContinuousSeries = chart.ContinuousSeries
type XAxis = chart.XAxis
type YAxis = chart.YAxis
type Style = chart.Style
type Tick = chart.Tick
type Renderable = chart.Renderable
type AnnotationSeries = chart.AnnotationSeries
type Value2 = chart.Value2

////go:linkname LegendThin gitee.com/quant1x/pkg/chart.LegendThin
//func LegendThin(c *Chart, userDefaults ...Style) Renderable
//
////go:linkname PNG gitee.com/quant1x/pkg/chart.PNG
//func PNG(width, height int) (chart.Renderer, error)

//go:linkname LastValueAnnotationSeries gitee.com/quant1x/pkg/chart.LastValueAnnotationSeries
func LastValueAnnotationSeries(innerSeries chart.ValuesProvider, vfs ...chart.ValueFormatter) chart.AnnotationSeries

// CreateChart 创建一个默认的图标
func CreateChart() chart.Chart {
	font, _ := GetDefaultFont()
	font.Bounds(1)
	lineChartStyle := chart.Style{
		Padding: chart.Box{
			Top: 20,
		},
	}
	graph := chart.Chart{
		Font:       font,
		Background: lineChartStyle,
	}
	return graph
}

func AddSeries(graph chart.Chart, series ...chart.Series) chart.Chart {
	graph.Series = append(graph.Series, series...)
	return graph
}

func Render(graph chart.Chart, code string) {
	graph.Elements = []chart.Renderable{chart.LegendThin(&graph)}
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.WriteFile(code+".png", buffer.Bytes(), api.CACHE_FILE_MODE)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Chart 图表
type Chart struct {
	chart.Chart
}

// NewChart 创建默认的图表
func NewChart() *Chart {
	font, _ := GetDefaultFont()
	font.Bounds(1)
	lineChartStyle := chart.Style{
		Padding: chart.Box{
			Top: 20,
		},
	}
	graph := chart.Chart{
		Font:       font,
		Background: lineChartStyle,
	}
	return &Chart{Chart: graph}
}

// AddSeries 添加图表序列
func (this *Chart) AddSeries(series ...chart.Series) {
	this.Series = append(this.Series, series...)
}

// Output 输出图表
func (this *Chart) Output(name string) error {
	this.Elements = []chart.Renderable{chart.LegendThin(&this.Chart)}
	buffer := bytes.NewBuffer([]byte{})
	err := this.Render(chart.PNG, buffer)
	if err != nil {
		return err
	}
	err = os.WriteFile(name+".png", buffer.Bytes(), api.CACHE_FILE_MODE)
	if err != nil {
		return err
	}
	return nil
}
