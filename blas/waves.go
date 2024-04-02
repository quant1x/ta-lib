package blas

import (
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/num"
)

// Wave 波浪 波峰波谷
type Wave struct {
	//data        []float64       // 原始数据
	diff        []int           // 一阶差分
	peaks       []num.DataPoint // 波峰位置存储
	valleys     []num.DataPoint // 波谷位置存储
	peakCount   int             // 所识别的波峰计数
	valleyCount int             // 所识别的波谷计数
}

// Waves 波浪 波峰波谷
type Waves struct {
	Data        []float64       // 原始数据
	Peaks       []num.DataPoint // 波峰位置存储
	Valleys     []num.DataPoint // 波谷位置存储
	PeakCount   int             // 所识别的波峰计数
	ValleyCount int             // 所识别的波谷计数
	Digits      int             // 小数点几位
}

//// 初始化波峰波谷
//func initializeWave(sample DataSample) *Wave {
//	n := sample.Len()
//	if n == 0 {
//		return nil
//	}
//	wave := Wave{
//		data:    make([]float64, n),
//		diff:    make([]int, n),
//		peaks:   make([]int, n),
//		valleys: make([]int, n),
//	}
//	for i := 0; i < n; i++ {
//		wave.data[i] = sample.Current(i)
//		wave.diff[i] = 0
//		wave.peaks[i] = -1
//		wave.valleys[i] = -1
//	}
//	wave.peakCount = 0
//	wave.valleyCount = 0
//	//step 1: 首先进行前向差分，并归一化
//	for i := 0; i < n-1; i++ {
//		a := wave.data[i]
//		b := wave.data[i+1]
//		sampleDiff := b - a
//		if sampleDiff > 0 {
//			wave.diff[i] = 1
//		} else if sampleDiff < 0 {
//			wave.diff[i] = -1
//		} else {
//			wave.diff[i] = 0
//		}
//	}
//
//	// step 2: 对相邻相等的点进行领边坡度处理
//	for i := 0; i < n-1; i++ {
//		if wave.diff[i] == 0 {
//			if i == (n - 2) {
//				if wave.diff[i-1] >= 0 {
//					wave.diff[i] = 1
//				} else {
//					wave.diff[i] = -1
//				}
//			} else {
//				if wave.diff[i+1] >= 0 {
//					wave.diff[i] = 1
//				} else {
//					wave.diff[i] = -1
//				}
//			}
//		}
//	}
//	// step 3: 对相邻相等的点进行领边坡度处理
//	for i := 0; i < n-1; i++ {
//		sampleDiff := wave.diff[i+1] - wave.diff[i]
//		if sampleDiff == -2 {
//			// 波峰识别
//			wave.peaks[wave.peakCount] = i + 1
//			wave.peakCount++
//		} else if sampleDiff == 2 {
//			// 波谷识别
//			wave.valleys[wave.valleyCount] = i + 1
//			wave.valleyCount++
//		}
//	}
//	wave.peaks = wave.peaks[:wave.peakCount]
//	wave.valleys = wave.valleys[:wave.valleyCount]
//	return &wave
//}

func findPeaks(n int, data func(x int) float64) Wave {
	wave := Wave{
		//data:    make([]float64, n),
		diff:    make([]int, n),
		peaks:   make([]num.DataPoint, n),
		valleys: make([]num.DataPoint, n),
	}
	for i := 0; i < n; i++ {
		wave.diff[i] = 0
		wave.peaks[i].X = -1
		wave.valleys[i].X = -1
	}
	wave.peakCount = 0
	wave.valleyCount = 0

	//step 1: 首先进行前向差分，并归一化
	for i := 0; i < n-1; i++ {
		a := data(i)
		b := data(i + 1)
		sampleDiff := b - a
		if sampleDiff > 0 {
			wave.diff[i] = 1
		} else if sampleDiff < 0 {
			wave.diff[i] = -1
		} else {
			wave.diff[i] = 0
		}
	}

	// step 2: 对相邻相等的点进行领边坡度处理
	for i := 0; i < n-1; i++ {
		if wave.diff[i] == 0 {
			if i == (n - 2) {
				if wave.diff[i-1] >= 0 {
					wave.diff[i] = 1
				} else {
					wave.diff[i] = -1
				}
			} else {
				if wave.diff[i+1] >= 0 {
					wave.diff[i] = 1
				} else {
					wave.diff[i] = -1
				}
			}
		}
	}
	// step 3: 对相邻相等的点进行领边坡度处理
	for i := 0; i < n-1; i++ {
		sampleDiff := wave.diff[i+1] - wave.diff[i]
		pos := i + 1
		if sampleDiff == -2 {
			// 波峰识别
			wave.peaks[wave.peakCount] = num.DataPoint{X: pos, Y: data(pos)}
			wave.peakCount++
		} else if sampleDiff == 2 {
			// 波谷识别
			wave.valleys[wave.valleyCount] = num.DataPoint{X: pos, Y: data(pos)}
			wave.valleyCount++
		}
	}
	wave.peaks = wave.peaks[:wave.peakCount]
	wave.valleys = wave.valleys[:wave.valleyCount]
	return wave
}

func v1FindPeaks(sample DataSample) Wave {
	n := sample.Len()
	wave := Wave{
		//data:    make([]float64, n),
		diff:    make([]int, n),
		peaks:   make([]num.DataPoint, n),
		valleys: make([]num.DataPoint, n),
	}
	for i := 0; i < n; i++ {
		wave.diff[i] = 0
		wave.peaks[i].X = -1
		wave.valleys[i].X = -1
	}
	wave.peakCount = 0
	wave.valleyCount = 0

	//step 1: 首先进行前向差分，并归一化
	for i := 0; i < n-1; i++ {
		a := sample.Current(i)
		b := sample.Current(i + 1)
		sampleDiff := b - a
		if sampleDiff > 0 {
			wave.diff[i] = 1
		} else if sampleDiff < 0 {
			wave.diff[i] = -1
		} else {
			wave.diff[i] = 0
		}
	}

	// step 2: 对相邻相等的点进行领边坡度处理
	for i := 0; i < n-1; i++ {
		if wave.diff[i] == 0 {
			if i == (n - 2) {
				if wave.diff[i-1] >= 0 {
					wave.diff[i] = 1
				} else {
					wave.diff[i] = -1
				}
			} else {
				if wave.diff[i+1] >= 0 {
					wave.diff[i] = 1
				} else {
					wave.diff[i] = -1
				}
			}
		}
	}
	// step 3: 对相邻相等的点进行领边坡度处理
	for i := 0; i < n-1; i++ {
		sampleDiff := wave.diff[i+1] - wave.diff[i]
		pos := i + 1
		if sampleDiff == -2 {
			// 波峰识别
			//if pos == 80 {
			//	fmt.Println("debug")
			//}
			if pos+1 < n && sample.High(pos) < sample.High(pos+1) {
				pos = pos + 1
			}
			price := sample.High(pos)
			wave.peaks[wave.peakCount] = num.DataPoint{X: pos, Y: price}
			wave.peakCount++
		} else if sampleDiff == 2 {
			// 波谷识别
			if pos+1 < n && sample.Low(pos) < sample.Low(pos+1) {
				pos = pos + 1
			}
			price := sample.Low(pos)
			wave.valleys[wave.valleyCount] = num.DataPoint{X: pos, Y: price}
			wave.valleyCount++
		}
	}
	wave.peaks = wave.peaks[:wave.peakCount]
	wave.valleys = wave.valleys[:wave.valleyCount]
	return wave
}

// NewWaves 创建一个新的波浪
func v1NewWaves(sample DataSample, code ...string) Waves {
	n := sample.Len()
	waves := Waves{
		Data: make([]float64, n),
	}
	for i := 0; i < n; i++ {
		waves.Data[i] = sample.Current(i)
	}
	digits := 2
	if len(code) > 0 {
		securityCode := exchange.CorrectSecurityCode(code[0])
		if info, ok := securities.CheckoutSecurityInfo(securityCode); ok {
			digits = int(info.DecimalPoint)
		}
	}
	waves.Digits = digits
	// 第一步, 高点
	highs := findPeaks(n, sample.High)
	waves.Peaks = highs.peaks
	waves.PeakCount = highs.peakCount
	// 第二步, 低点
	lows := findPeaks(n, sample.Low)
	waves.Valleys = lows.valleys
	waves.ValleyCount = lows.valleyCount
	return waves
}

// NewWaves 创建一个新的波浪
func NewWaves(sample DataSample, code ...string) Waves {
	n := sample.Len()
	waves := Waves{
		Data: make([]float64, n),
	}
	for i := 0; i < n; i++ {
		waves.Data[i] = sample.Current(i)
	}
	digits := 2
	if len(code) > 0 {
		securityCode := exchange.CorrectSecurityCode(code[0])
		if info, ok := securities.CheckoutSecurityInfo(securityCode); ok {
			digits = int(info.DecimalPoint)
		}
	}
	waves.Digits = digits
	wave := v1FindPeaks(sample)
	waves.Peaks = wave.peaks
	waves.PeakCount = wave.peakCount
	waves.Valleys = wave.valleys
	waves.ValleyCount = wave.valleyCount
	return waves
}
