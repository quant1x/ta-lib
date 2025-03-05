package algorithms

import (
	"math"
	"testing"
)

func TestQuantile(t *testing.T) {
	tests := []struct {
		data []float64
		p    float64
		want float64
	}{
		{[]float64{1, 2, 3, 4, 5}, 0.9, 4.6}, // 线性插值
		{[]float64{1, 2, 3, 4}, 0.75, 3.25},  // 四分位数
		{[]float64{10}, 0.5, 10},             // 单元素
		{[]float64{5, 1, 3, 2, 4}, 0.5, 3},   // 乱序输入
	}

	for _, tt := range tests {
		got, _ := Quantile(tt.data, tt.p)
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("got %.2f, want %.2f", got, tt.want)
		}
	}
}
