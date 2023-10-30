package indicators

import (
	"fmt"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pandas/stat"
	"gitee.com/quant1x/ta-lib/linear"
)

// M M头跌破颈线
func M(raw pandas.DataFrame, argv ...bool) (stat.DType, bool) {
	const (
		__delta = float64(0)
		//__delta = linear.TrendDelta
	)
	var (
		__debug        bool = false
		__ignoreSignal      = false
	)
	if raw.Nrow() < linear.MaximumTrendPeriod {
		return 0.00, false
	}
	if len(argv) > 0 {
		__debug = argv[0]
	}
	if len(argv) > 1 {
		__ignoreSignal = argv[1]
	}
	df := raw.Subset(raw.Nrow()-linear.MaximumTrendPeriod, raw.Nrow())
	//fmt.Println(df)
	// 获取全部的波峰的数据
	vh := df.Col("high")
	v := vh.DTypes()
	_, _, iMax, vMax := linear.PeakDetect(v[:], __delta)
	if len(iMax) < 1 {
		return 0.00, false
	}
	// 获取全部的波谷的数据
	vl := df.Col("close")
	v = vl.DTypes()
	iMin, vMin, _, _ := linear.PeakDetect(v[:], __delta)
	var ylX1, ylX2 int
	var ylY1, ylY2 float64
	//var ylSlope float64
	if len(iMax) >= 2 {
		w := len(iMax)
		ylX1 = iMax[w-2]
		ylY1 = vMax[w-2]
		ylX2 = iMax[w-1]
		ylY2 = vMax[w-1]
	} else {
		return 0.00, false
	}
	// 找到最近的一个波谷
	minX := -1
	minY := float64(-1) // 颈线
	for i := 0; i < len(iMin); i++ {
		if iMin[i] > ylX1 && iMin[i] < ylX2 {
			minX = iMin[i]
			minY = vMin[i]
			break
		}
	}

	if minX > 0 {
		dbH, dbL := MAX_GO(ylY1, ylY2)
		pLow := minY*2 - dbH
		pHigh := minY*2 - dbL
		dates := df.Col("date").Strings()
		closes := df.Col("close").DTypes()
		cl := len(closes)
		if __ignoreSignal || (closes[cl-2] < closes[cl-1] && closes[cl-2] < ylY1 && closes[cl-1] > ylY1) {
			if __debug {
				fmt.Printf("       M头: 左高(%s)=%f, 右高(%s)=%f, 颈线(%s)=%f\n", dates[ylX1], ylY1, dates[ylX2], ylY2, dates[minX], minY)
				fmt.Printf("M头跌破颈线: 颈线(%s)=%f, %s\n", dates[minX], minY, dates[cl-1])
				fmt.Printf("预计目标位置: 最低=%f, 最高=%f\n", pLow, pHigh)
			}
			return pHigh, true
		}
	}

	return 0.00, false
}
