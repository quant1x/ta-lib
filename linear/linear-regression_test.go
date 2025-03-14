package linear

import (
	"fmt"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/num"
	"testing"
)

func TestPredictStock(t *testing.T) {
	df := factors.KLine("002528")
	fmt.Println(df)
	length := df.Nrow() - 1
	df1 := df.Subset(length-3, length)
	fmt.Println(df1)
	CLOSE := df1.Col("low").DTypes()
	//CLOSE = []float64{1, 2, 3, 4, 5}
	data_len := len(CLOSE)
	fmt.Printf("raw   data length: %d \n", data_len)
	// 去掉最后1天的数据
	y := CLOSE[:data_len]
	y_length := len(y)
	fmt.Printf("train data length: %d, last data[%d]=%f \n", y_length, (y_length - 1), y[y_length-1])
	x := make([]float64, len(y))
	for i, v := range y {
		x[i] = float64(i)
		_ = v
	}

	fmt.Println("------------------------------------------------------------")
	p1 := num.PolyFit(y, x, 2)
	fmt.Println("p1 =", p1)
	fmt.Println("------------------------------------------------------------")

	k, b := LeastSquares(x, y)
	// 预测最后1天的下一个交易日的数据
	no := y_length
	fmt.Printf("no: %d, predicting...\n", no)
	p := Predict(float64(no), k, b)
	fmt.Printf("no: %d, predicted=%f\n", no, p)
}

func TestPolyFit(t *testing.T) {
	x := []float64{0.0, 0.1, 0.2, 0.3, 0.5, 0.8, 1.0}
	y := []float64{1.0, 0.41, 0.50, 0.61, 0.91, 2.02, 2.46}
	A := num.PolyFit(x, y, 2)
	fmt.Println("A =", A)

	//A2 := []float64{3.131561350718812, -1.2400367769976413, 0.7355767301905694}
	z1 := num.PolyVal(A, x)
	fmt.Println("z1 =", z1)

	W := 5
	A2 := num.PolyFit(y, num.Range[float64](W), 1)
	fmt.Println("A2 =", A2)
	x2 := num.Repeat[float64](float64(W), W)
	z2 := num.PolyVal(A2, x2)
	fmt.Println("z2 =", z2)
}
