package indicators

// FeatureSar SAR指标特征数据结构
type FeatureSar struct {
	Pos    int     // 坐标位置
	Bull   bool    // 当前多空
	Af     float64 // 加速因子(Acceleration Factor)
	Ep     float64 // 极值点(Extreme Point)
	Sar    float64 // SAR[Pos]
	High   float64 // pos周期最高价
	Low    float64 // pos周期最低价
	Period int     // 周期数, 上涨趋势, 周期数大于0, 下跌趋势, 周期数小于0, 绝对值就是已过多少天
}

const (
	// SarAccelerationFactor SAR 加速因子
	SarAccelerationFactor = 0.02
	// SarAccelerationFactorLimit SAR加速因子最大值
	SarAccelerationFactorLimit = 0.20
)

// SAR 停损转向操作点指标
func SAR(highs, lows []float64) []FeatureSar {
	return StopAndReverse(true, highs, lows, SarAccelerationFactor, SarAccelerationFactorLimit)
}

// StopAndReverse 停损转向操作点指标
//
//	SAR指标又叫抛物线指标或停损转向操作点指标, 其全称叫“Stop and Reverse, 缩写SAR”,
//	是由美国技术分析大师威尔斯-威尔德(Wells Wilder)所创造的, 是一种简单易学,比较准确的中短期技术分析工具.
//	https://baike.baidu.com/item/SAR%E6%8C%87%E6%A0%87
func StopAndReverse(firstIsBull bool, highs, lows []float64, accelerationFactor, accelerationFactorLimit float64) []FeatureSar {
	return v2Sar(firstIsBull, highs, lows, accelerationFactor, accelerationFactorLimit)
}

func v1Sar(firstIsBull bool, highs, lows []float64, accelerationFactor, accelerationFactorLimit float64) []FeatureSar {
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

func v2Sar(firstIsBull bool, highs, lows []float64, accelerationFactor, accelerationFactorLimit float64) []FeatureSar {
	length := len(highs)
	data := make([]FeatureSar, length)

	// 第一个bar
	data[0].Pos = 0
	data[0].Bull = firstIsBull
	// 加速因子
	data[0].Af = accelerationFactor
	// 极值点
	//extremePoint := highs[0]
	data[0].Ep = highs[0]
	data[0].Sar = lows[0]
	data[0].High = highs[0]
	data[0].Low = lows[0]
	data[0].Period = 1
	for i := 1; i < length; i++ {
		//data[i] = sarIncr(data[i-1], accelerationFactor, accelerationFactorLimit, highs[i], lows[i])
		data[i] = data[i-1].RawIncr(accelerationFactor, accelerationFactorLimit, highs[i], lows[i])
	}
	return data
}

// 增量计算
func sarIncr(last FeatureSar, accelerationFactor, accelerationFactorLimit float64, high, low float64) FeatureSar {
	current := last
	current.Pos++
	current.High = high
	current.Low = low
	// 1. 初次赋值
	if last.Bull {
		// 多头
		if high > last.Ep {
			// 创新高
			current.Ep = high
			current.Af = min(last.Af+accelerationFactor, accelerationFactorLimit)
		}
	} else {
		// 空头
		if low < last.Ep {
			// 创新低
			current.Ep = low
			current.Af = min(last.Af+accelerationFactor, accelerationFactorLimit)
		}
	}
	// 2. 计算SAR
	current.Sar = last.Sar + current.Af*(current.Ep-last.Sar)
	// 3. 修正SAR
	if last.Bull {
		current.Sar = max(last.Sar, min(current.Sar, low, last.Low))
	} else {
		current.Sar = min(last.Sar, max(current.Sar, high, last.High))
	}
	// 4. 判断多空
	if last.Bull {
		// 多
		if low < current.Sar {
			// 向下跌破, 转空
			current.Bull = false
			current.Ep = low
			current.Af = accelerationFactor
			if last.High == last.Ep {
				// 紧邻即高点
				current.Sar = last.Ep
			} else {
				current.Sar = last.Ep + current.Af*(current.Ep-last.Ep)
			}
		}
	} else {
		// 空
		if high > current.Sar {
			// 向上突破, 转多
			current.Bull = true
			current.Ep = high
			current.Af = accelerationFactor
			current.Sar = min(low, last.Low)
		}
	}
	return current
}

// Incr 增量计算
func (s FeatureSar) Incr(high, low float64) FeatureSar {
	return s.RawIncr(SarAccelerationFactor, SarAccelerationFactorLimit, high, low)
}

// RawIncr 增加1天的数据
func (s FeatureSar) RawIncr(accelerationFactor, accelerationFactorLimit float64, high, low float64) FeatureSar {
	current := s
	current.Pos++
	current.High = high
	current.Low = low
	// 1. 初次赋值
	if s.Bull {
		// 多头
		if high > s.Ep {
			// 创新高
			current.Ep = high
			current.Af = min(s.Af+accelerationFactor, accelerationFactorLimit)
		}
	} else {
		// 空头
		if low < s.Ep {
			// 创新低
			current.Ep = low
			current.Af = min(s.Af+accelerationFactor, accelerationFactorLimit)
		}
	}
	// 2. 计算SAR
	current.Sar = s.Sar + current.Af*(current.Ep-s.Sar)
	// 3. 修正SAR
	if s.Bull {
		current.Sar = max(s.Sar, min(current.Sar, low, s.Low))
	} else {
		current.Sar = min(s.Sar, max(current.Sar, high, s.High))
	}
	// 4. 判断多空
	if s.Bull {
		// 多
		if low < current.Sar {
			// 向下跌破, 转空
			current.Bull = false
			current.Ep = low
			current.Af = accelerationFactor
			if s.High == s.Ep {
				// 紧邻即高点
				current.Sar = s.Ep
			} else {
				current.Sar = s.Ep + current.Af*(current.Ep-s.Ep)
			}
			current.Period = -1
		} else {
			current.Period++
		}
	} else {
		// 空
		if high > current.Sar {
			// 向上突破, 转多
			current.Bull = true
			current.Ep = high
			current.Af = accelerationFactor
			current.Sar = min(low, s.Low)
			current.Period = 1
		} else {
			current.Period--
		}
	}
	return current
}
