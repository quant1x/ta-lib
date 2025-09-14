package blas

import (
	"github.com/quant1x/exchange"
	"github.com/quant1x/gotdx/securities"
	"github.com/quant1x/num"
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
	State       int8            // 点状态
	LastHigh    num.DataPoint   // 最新的高点
	LastLow     num.DataPoint   // 最新的地点
}

func (this Waves) Len() int {
	return len(this.Data)
}

// FindPeaks 波峰波谷
//
//	搜寻收盘价的高低点
func FindPeaks(n int, data func(x int) float64) Wave {
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

// PeaksAndValleys 创建一个新的波浪
//
//	高点的波峰, 低点的波谷
func PeaksAndValleys(sample DataSample, code ...string) Waves {
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
	wave := searchHighsAndLows(sample)
	waves.Peaks = wave.peaks
	waves.PeakCount = wave.peakCount
	waves.Valleys = wave.valleys
	waves.ValleyCount = wave.valleyCount
	return waves
}

// 搜寻高点和低点
func searchHighsAndLows(sample DataSample) Wave {
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
			//if i+2 < n {
			//	tmp := []float64{sample.High(i), sample.High(i + 1), sample.High(i + 2)}
			//	n := num.ArgMax(tmp)
			//	pos = i + n
			//}
			price := sample.High(pos)
			wave.peaks[wave.peakCount] = num.DataPoint{X: pos, Y: price}
			wave.peakCount++
		} else if sampleDiff == 2 {
			// 波谷识别
			//if i+2 < n {
			//	tmp := []float64{sample.Low(i), sample.Low(i + 1), sample.Low(i + 2)}
			//	n := num.ArgMin(tmp)
			//	pos = i + n
			//}
			price := sample.Low(pos)
			wave.valleys[wave.valleyCount] = num.DataPoint{X: pos, Y: price}
			wave.valleyCount++
		}
	}
	wave.peaks = wave.peaks[:wave.peakCount]
	wave.valleys = wave.valleys[:wave.valleyCount]
	return wave
}

// HighAndLow 关注低点高点变化的波浪
func HighAndLow(sample DataSample, code ...string) Waves {
	n := sample.Len()
	waves := Waves{
		Data:    make([]float64, n),
		Peaks:   make([]num.DataPoint, n),
		Valleys: make([]num.DataPoint, n),
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
	waves.PeakCount = 0
	waves.ValleyCount = 0
	waves.State = 0
	pos := 0
	waves.LastHigh = num.DataPoint{X: pos, Y: sample.High(pos)}
	waves.LastLow = num.DataPoint{X: pos, Y: sample.Low(pos)}
	for i := 1; i < n; i++ {
		lastHigh := num.DataPoint{X: i, Y: sample.High(i)}
		lastLow := num.DataPoint{X: i, Y: sample.Low(i)}
		diffHigh := lastHigh.Y - waves.LastHigh.Y
		diffLow := lastLow.Y - waves.LastLow.Y
		switch waves.State {
		case 0: // 初始状态
			if diffHigh > 0 {
				waves.State = 1
			} else if diffLow < 0 {
				waves.State = -1
			}
		case 1: // 升高
			if diffHigh <= 0 {
				waves.State = -1
				waves.Peaks[waves.PeakCount] = waves.LastHigh
				waves.PeakCount++
			}
		case -1: // 降低
			if diffLow >= 0 {
				waves.State = 1
				waves.Valleys[waves.ValleyCount] = waves.LastLow
				waves.ValleyCount++
			}
		}
		waves.LastHigh = lastHigh
		waves.LastLow = lastLow
	}
	waves.Peaks = waves.Peaks[:waves.PeakCount]
	waves.Valleys = waves.Valleys[:waves.ValleyCount]
	return waves
}
