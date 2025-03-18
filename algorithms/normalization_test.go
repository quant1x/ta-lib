package algorithms

import (
	"fmt"
	"slices"
	"testing"
)

func TestNormalize(t *testing.T) {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	slices.Max(data)
	// 使用最小-最大归一化
	normalized, err := Normalize(data, MethodMinMax)
	if err != nil {
		panic(err)
	}
	fmt.Println("Min-Max:", normalized)

	// 使用Z-Score标准化
	normalized, err = Normalize(data, MethodZScore)
	if err != nil {
		panic(err)
	}
	fmt.Println("Z-Score:", normalized)

	// 鲁棒标准化
	normalized, err = Normalize(data, MethodRobust)
	if err != nil {
		panic(err)
	}
	fmt.Println("Robust:", normalized)

	normalized, err = Normalize(data, MethodL2)
	fmt.Println("L2 Norm:", normalized)
}
