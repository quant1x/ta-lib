package labs

import (
	"fmt"
	"testing"
)

func TestBoolRelExtrema(t *testing.T) {
	// 示例测试
	data := []float64{1, 2, 3, 2, 1}
	comparator := func(a, b float64) bool { return a > b }
	order := 1
	mode := "clip"

	extrema, err := ArgRelExtrema(data, comparator, order, mode)
	if err != nil {
		fmt.Println("错误:", err)
		return
	}

	fmt.Println(extrema) // 输出: [false false true false false]
}
