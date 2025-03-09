package algorithms

import (
	"fmt"
	"testing"
)

// 示例：检测RSI趋势（类似Python代码逻辑）
func checkRSIDivergence(rsi6, rsi12, rsi24 []float64, trendLength int) string {
	// 取最后trendLength个数据点
	x := make([]float64, trendLength)
	for i := 0; i < trendLength; i++ {
		x[i] = float64(i)
	}

	// 拟合线性趋势
	coef6 := Polyfit(x, rsi6[len(rsi6)-trendLength:], 1)
	coef12 := Polyfit(x, rsi12[len(rsi12)-trendLength:], 1)
	coef24 := Polyfit(x, rsi24[len(rsi24)-trendLength:], 1)

	// 判断趋势方向
	if coef6[0] > 0 && coef12[0] > 0 && coef24[0] > 0 {
		return "多头排列"
	} else if coef6[0] < 0 && coef12[0] < 0 && coef24[0] < 0 {
		return "空头排列"
	}
	return "无显著趋势"
}

func TestPolyfit(t *testing.T) {
	// 示例数据（RSI6/12/24的最后3个值）
	rsi6 := []float64{60, 65, 68} // 上升趋势
	rsi12 := []float64{55, 62, 67}
	rsi24 := []float64{50, 58, 65}

	result := checkRSIDivergence(rsi6, rsi12, rsi24, 3)
	fmt.Println("趋势判断:", result) // 应输出"多头排列"
}
