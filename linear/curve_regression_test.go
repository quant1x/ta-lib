package linear

import (
	"fmt"
	"gitee.com/quant1x/ta-lib/testfiles"
	"testing"
)

func TestCurveRegression(t *testing.T) {
	df := testfiles.LoadTestData()
	df = df.Subset(0, df.Nrow()-1)
	fmt.Println(df)
	N := 3
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
