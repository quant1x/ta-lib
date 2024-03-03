package plot

import (
	"bytes"
	"fmt"
	"gitee.com/quant1x/gox/api"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
	"os"
)

const (
	DotWidth = 3 // 点的宽度
)

var (
	ColorRed   = drawing.Color{R: 0xff, G: 0x00, B: 0x00, A: 0xff} // 红色
	ColorGreen = chart.ColorGreen                                  // 绿色
	ColorBlue  = chart.ColorBlue                                   // 蓝色
)

var (
	DashedLine = []float64{5.0, 5.0} // 虚线
)

func CreateChart() chart.Chart {
	font, _ := GetDefaultFont()
	font.Bounds(1)
	lineChartStyle := chart.Style{
		Padding: chart.Box{
			//Top:  10,
			Left: 50,
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
