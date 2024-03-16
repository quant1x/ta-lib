package linear

import (
	"fmt"
	"gitee.com/quant1x/engine/factors"
	"testing"
)

func TestCurveRegression(t *testing.T) {
	code := "688351"
	code = "sh000001"
	df := factors.KLine(code)
	df = df.Subset(0, df.Nrow()-1)
	fmt.Println(df)
	N := 5
	V := df.Col("open")
	d := CurveRegression(V, N)
	fmt.Println(d)
	V = df.Col("close")
	d = CurveRegression(V, N)
	fmt.Println(d)
	V = df.Col("high")
	d = CurveRegression(V, N)
	fmt.Println(d)
	V = df.Col("low")
	d = CurveRegression(V, N)
	fmt.Println(d)
}
