package linear

import (
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/pandas"
	"testing"
)

func TestCurveRegression(t *testing.T) {
	code := "600839"
	code = exchange.CorrectSecurityCode(code)
	date := exchange.GetCurrentlyDay()
	rawData := base.CheckoutKLines(code, date)
	df := pandas.LoadStructs(rawData)
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
