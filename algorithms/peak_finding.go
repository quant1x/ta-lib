package algorithms

import (
	"errors"
)

// adjustIndex 处理索引越界，根据mode返回有效索引
func adjustIndex(i, length int, mode string) int {
	if length == 0 {
		return 0
	}
	if i < 0 {
		switch mode {
		case "clip":
			return 0
		case "wrap":
			return (i%length + length) % length
		default:
			return 0 // 默认为clip模式
		}
	} else if i >= length {
		switch mode {
		case "clip":
			return length - 1
		case "wrap":
			return i % length
		default:
			return length - 1 // 默认为clip模式
		}
	}
	return i
}

// ArgRelExtrema 检测数据中的相对极值点
func ArgRelExtrema(data []float64, comparator func(a, b float64) bool, order int, mode string) ([]bool, error) {
	if order < 1 {
		return nil, errors.New("order必须为大于等于1的整数")
	}

	dataLen := len(data)
	results := make([]bool, dataLen)
	for i := range results {
		results[i] = true
	}

	for shift := 1; shift <= order; shift++ {
		anyTrue := false
		for i := 0; i < dataLen; i++ {
			// 计算左右偏移后的索引
			plusIndex := adjustIndex(i+shift, dataLen, mode)
			minusIndex := adjustIndex(i-shift, dataLen, mode)

			// 比较当前值与左右值
			current := data[i]
			plusVal := data[plusIndex]
			minusVal := data[minusIndex]

			// 应用比较函数
			cmpPlus := comparator(current, plusVal)
			cmpMinus := comparator(current, minusVal)

			// 更新结果数组
			results[i] = results[i] && cmpPlus && cmpMinus

			if results[i] {
				anyTrue = true
			}
		}

		// 如果所有结果都为false，提前终止
		if !anyTrue {
			break
		}
	}

	return results, nil
}
