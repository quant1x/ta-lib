package indicators

type FeatureSar struct {
	Pos  int     // 坐标位置
	Bull bool    // 当前多空
	Af   float64 // 加速因子(Acceleration Factor)
	Ep   float64 // 极值点(Extreme Point)
	Sar  float64 // SAR[Pos]
}

// SAR指标又叫抛物线指标或停损转向操作点指标, 其全称叫“Stop and Reverse, 缩写SAR”,
// 是由美国技术分析大师威尔斯-威尔德(Wells Wilder)所创造的, 是一种简单易学,比较准确的中短期技术分析工具.
// https://baike.baidu.com/item/SAR%E6%8C%87%E6%A0%87
func tdx_sar(firstIsBull bool, highs, lows []float64, accelerationFactor, accelerationFactorLimit float64) []FeatureSar {
	length := len(highs)
	data := make([]FeatureSar, length)

	// 第一个bar
	data[0].Pos = 0
	data[0].Bull = firstIsBull
	// 加速因子
	data[0].Af = accelerationFactor
	// 极值点
	//extremePoint := highs[0]
	data[0].Ep = highs[data[0].Pos]
	data[0].Sar = lows[data[0].Pos]
	for i := 1; i < length; i++ {
		prevSar := data[i-1]
		data[i] = prevSar
		current := &data[i]
		current.Pos = i
		current.Bull = prevSar.Bull
		// 1. 初次赋值
		if data[i-1].Bull {
			// 多头
			if highs[i] > prevSar.Ep {
				// 创新高
				current.Ep = highs[i]
				current.Af = min(prevSar.Af+accelerationFactor, accelerationFactorLimit)
			}
		} else {
			// 空头
			if lows[i] < prevSar.Ep {
				// 创新低
				current.Ep = lows[i]
				current.Af = min(prevSar.Af+accelerationFactor, accelerationFactorLimit)
			}
		}
		// 2. 计算SAR
		current.Sar = prevSar.Sar + current.Af*(current.Ep-prevSar.Sar)
		// 3. 修正SAR
		if prevSar.Bull {
			current.Sar = max(prevSar.Sar, min(current.Sar, lows[i], lows[i-1]))
		} else {
			current.Sar = min(prevSar.Sar, max(current.Sar, highs[i], highs[i-1]))
		}
		// 4. 判断多空
		if prevSar.Bull {
			// 多
			if lows[i] < current.Sar {
				// 向下跌破, 转空
				current.Bull = false
				current.Ep = lows[i]
				current.Af = accelerationFactor
				if highs[i-1] == prevSar.Ep {
					// 紧邻即高点
					current.Sar = prevSar.Ep
				} else {
					current.Sar = prevSar.Ep + current.Af*(current.Ep-prevSar.Ep)
				}
			}
		} else {
			// 空
			if highs[i] > current.Sar {
				// 向上突破, 转多
				current.Bull = true
				current.Ep = highs[i]
				current.Af = accelerationFactor
				current.Sar = min(lows[i], lows[i-1])
			}
		}
	}
	return data
}
