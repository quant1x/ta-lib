package algorithms

import (
	"fmt"
	"slices"
	"sort"
)

// Quantile 计算分位数 (0 <= p <= 1)
func Quantile(data []float64, p float64) (float64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("empty dataset")
	}
	if p < 0 || p > 1 {
		return 0, fmt.Errorf("p must be between 0 and 1")
	}

	// 创建数据副本避免修改原数据
	//sorted := make([]float64, len(data))
	//copy(sorted, data)
	sorted := slices.Clone(data)
	sort.Float64s(sorted)

	n := float64(len(sorted))
	pos := (n - 1) * p // 位置计算公式

	// 获取前后索引
	lower := int(pos)
	upper := lower + 1
	weight := pos - float64(lower)

	// 处理边界情况
	if upper >= len(sorted) {
		return sorted[len(sorted)-1], nil
	}

	// 线性插值
	return sorted[lower]*(1-weight) + sorted[upper]*weight, nil
}
