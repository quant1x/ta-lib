package algorithms

import (
	"errors"
	"math"
	"slices"
	"sort"
)

// 归一化方法常量定义
const (
	MethodMinMax          = "minmax"           // 最小-最大归一化（线性缩放）
	MethodZScore          = "zscore"           // Z-Score标准化
	MethodRobust          = "robust"           // 鲁棒标准化
	MethodL2              = "l2"               // L2正则化（单位向量）
	MethodLog             = "log"              // 对数变换
	MethodQuantileUniform = "quantile_uniform" // 映射到均匀分布
	MethodQuantileNormal  = "quantile_normal"  // 映射到正态分布
)

// 自定义错误类型
var (
	ErrEmptyData        = errors.New("data cannot be empty")
	ErrInvalidMethod    = errors.New("invalid normalization method")
	ErrConstantFeature  = errors.New("feature has constant value (zero variance)")
	ErrInsufficientData = errors.New("insufficient data points for quantile transformation")
)

// Normalize 归一化入口函数
// 方法选择建议：
//   - 异常值多 → 鲁棒标准化、分位数变换
//   - 稀疏数据 → L2正则化、绝对最大值归一化
//   - 深度学习 → 最小-最大（适配Sigmoid/ReLU）
//   - 正态假设 → Z-Score或分位数正态变换
func Normalize(data []float64, method string) ([]float64, error) {
	if len(data) == 0 {
		return nil, ErrEmptyData
	}

	switch method {
	case MethodMinMax:
		return minMaxScale(data), nil
	case MethodZScore:
		return zScoreScale(data)
	case MethodRobust:
		return robustScale(data)
	case MethodL2:
		return l2Normalize(data), nil
	case MethodLog:
		return logTransform(data), nil
	case MethodQuantileUniform:
		return quantileTransform(data, method)
	case MethodQuantileNormal:
		return quantileTransform(data, method)
	default:
		return nil, ErrInvalidMethod
	}
}

// region 线性归一化方法

// minMaxScale 最小-最大归一化（线性缩放）
//
//	公式：(x - dataMin) / (dataMax - dataMin)
func minMaxScale(data []float64) []float64 {
	dataMin := slices.Min(data)
	dataMax := slices.Max(data)
	if dataMin == dataMax {
		return data // 避免除零错误
	}

	result := make([]float64, len(data))
	for i, v := range data {
		result[i] = (v - dataMin) / (dataMax - dataMin)
	}
	return result
}

// endregion

// region 标准化方法

// zScoreScale Z-Score标准化
//
//	公式：(x - mean) / stdDev
func zScoreScale(data []float64) ([]float64, error) {
	m := mean(data)
	std := stdDev(data)
	if std == 0 {
		return nil, ErrConstantFeature
	}

	result := make([]float64, len(data))
	for i, v := range data {
		result[i] = (v - m) / std
	}
	return result, nil
}

// robustScale 鲁棒标准化
//
//	公式：(x - median) / IQR
func robustScale(data []float64) ([]float64, error) {
	sorted := slices.Clone(data)
	slices.Sort(sorted)

	q1 := getPercentile(sorted, 25)
	q3 := getPercentile(sorted, 75)
	iqr := q3 - q1
	median := getMedian(sorted)

	if iqr == 0 {
		return nil, ErrConstantFeature
	}

	result := make([]float64, len(data))
	for i, v := range data {
		result[i] = (v - median) / iqr
	}
	return result, nil
}

// endregion

// region 向量归一化方法

// l2Normalize L2正则化（单位向量）
//
//	公式：x / ||x||₂
func l2Normalize(data []float64) []float64 {
	norm := 0.0
	for _, v := range data {
		norm += v * v
	}
	norm = math.Sqrt(norm)
	if norm == 0 {
		return data
	}

	result := make([]float64, len(data))
	for i, v := range data {
		result[i] = v / norm
	}
	return result
}

// endregion

// region 非线性归一化方法

// logTransform 对数变换
//
//	公式：log(x/dataMax + ε)
func logTransform(data []float64) []float64 {
	dataMax := slices.Max(data)
	result := make([]float64, len(data))
	for i, v := range data {
		adjusted := (v / dataMax) + 1e-9 // 防止log(0)
		result[i] = math.Log(adjusted)
	}
	return result
}

// quantileTransform 分位数变换
//
//	支持均匀分布(0-1)和正态分布转换
func quantileTransform(data []float64, distribution string) ([]float64, error) {
	if len(data) < 2 {
		return nil, ErrInsufficientData
	}

	sorted := slices.Clone(data)
	slices.Sort(sorted)
	n := float64(len(sorted))
	result := make([]float64, len(data))

	for i, v := range data {
		pos := sort.Search(len(sorted), func(j int) bool {
			return sorted[j] >= v
		})
		fraction := float64(pos) / n

		switch distribution {
		case MethodQuantileUniform:
			result[i] = fraction
		case MethodQuantileNormal:
			result[i] = probit(fraction)
		default:
			return nil, ErrInvalidMethod
		}
	}
	return result, nil
}

// endregion

// region 辅助函数

// mean 计算均值
func mean(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

// stdDev 计算样本标准差
func stdDev(data []float64) float64 {
	m := mean(data)
	sum := 0.0
	for _, v := range data {
		sum += math.Pow(v-m, 2)
	}
	variance := sum / float64(len(data)-1)
	return math.Sqrt(variance)
}

// getMedian 计算中位数
func getMedian(sorted []float64) float64 {
	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}

// getPercentile 计算百分位数（线性插值法）
func getPercentile(sorted []float64, percent int) float64 {
	index := float64(percent) / 100 * float64(len(sorted)-1)
	lower := math.Floor(index)
	upper := math.Ceil(index)
	if lower == upper {
		return sorted[int(index)]
	}
	weight := index - lower
	return sorted[int(lower)]*(1-weight) + sorted[int(upper)]*weight
}

// probit 标准正态分布分位数函数（逆CDF）
func probit(p float64) float64 {
	if p <= 0 || p >= 1 {
		return math.NaN()
	}

	t := math.Sqrt(-2.0 * math.Log(math.Min(p, 1-p)))
	c := []float64{2.515517, 0.802853, 0.010328}
	d := []float64{1.432788, 0.189269, 0.001308}
	numerator := c[0] + c[1]*t + c[2]*t*t
	denominator := 1.0 + d[0]*t + d[1]*t*t + d[2]*t*t*t

	if p > 0.5 {
		return t - numerator/denominator
	}
	return -(t - numerator/denominator)
}

// endregion
