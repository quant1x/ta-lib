package algorithms

import (
	"errors"
	"math"
)

// VPVR 计算成交量加权的价格区间最大值
func VPVR(closePrices []float64, volumes []float64, bins int) (float64, error) {
	// 参数校验
	if len(closePrices) == 0 || len(volumes) == 0 {
		return 0, errors.New("empty input data")
	}
	if len(closePrices) != len(volumes) {
		return 0, errors.New("close and volume arrays must have same length")
	}
	if bins <= 0 {
		return 0, errors.New("bins must be positive integer")
	}

	// 1. 计算价格范围
	minPrice, maxPrice := findMinMax(closePrices)

	// 2. 生成分箱边界
	edges := generateEdges(minPrice, maxPrice, bins)

	// 3. 创建直方图容器
	hist := make([]float64, bins)

	// 4. 填充直方图
	binWidth := (maxPrice - minPrice) / float64(bins)
	for i, price := range closePrices {
		// 计算归属的bin索引
		binIndex := int(math.Floor((price - minPrice) / binWidth))

		// 处理最大值边界情况
		if binIndex == bins {
			binIndex = bins - 1
		}

		// 累加成交量
		if binIndex >= 0 && binIndex < bins {
			hist[binIndex] += volumes[i]
		}
	}

	// 5. 寻找最大成交量对应的bin
	maxIndex := 0
	maxValue := hist[0]
	for i, value := range hist {
		if value > maxValue {
			maxValue = value
			maxIndex = i
		}
	}

	// 返回对应区间的左边界
	return edges[maxIndex], nil
}

// 辅助函数：寻找最小最大值
func findMinMax(arr []float64) (min, max float64) {
	min = arr[0]
	max = arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return
}

// 辅助函数：生成分箱边界
func generateEdges(min, max float64, bins int) []float64 {
	edges := make([]float64, bins)
	binWidth := (max - min) / float64(bins)

	for i := range edges {
		edges[i] = min + float64(i)*binWidth
	}
	return edges
}
