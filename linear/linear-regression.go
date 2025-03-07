package linear

import (
	"gitee.com/quant1x/gox/logger"
)

// LeastSquares 最简单的最小二乘法
func LeastSquares(x []float64, y []float64) (slope float64, intercept float64) {
	// x是横坐标数据,y是纵坐标数据
	// a是斜率，b是截距
	xi := float64(0)
	x2 := float64(0)
	yi := float64(0)
	xy := float64(0)

	if len(x) != len(y) {
		logger.Debugf("最小二乘时，两数组长度不一致!")
		return
	} else {
		xLen := len(x)
		length := float64(xLen)
		window := 5
		if xLen <= window {
			window = xLen
		}
		for i := xLen - window; i < xLen; i++ {
			xi += x[i]
			x2 += x[i] * x[i]
			yi += y[i]
			xy += x[i] * y[i]
		}
		slope = (yi*xi - xy*length) / (xi*xi - x2*length)
		intercept = (yi*x2 - xy*xi) / (x2*length - xi*xi)
	}
	return
}

func Predict(y, slope, intercept float64) float64 {
	return y*slope + intercept
}
