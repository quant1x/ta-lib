package sample

import (
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
)

// ConfidenceInterval 置信区间(Confidence Interval, CI)
//
//	df必须要用一个字段为N, 作为时间序列
func ConfidenceInterval(df pandas.DataFrame, argv ...int) pandas.DataFrame {
	var (
		LEN   = df.Nrow()
		CLOSE = df.ColAsNDArray("close")
		HIGH  = df.ColAsNDArray("high")
		LOW   = df.ColAsNDArray("low")
		N     = df.ColAsNDArray("N").DTypes()
		CI    = 0.9500 // 95%的置信区间
	)
	if len(argv) > 0 {
		__n := argv[0]
		N = num.Repeat[num.DType](num.DType(__n), LEN)
	}
	mid := MA(CLOSE, N)
	variance := STD(CLOSE, N)
	Z := num.ConfidenceIntervalToZscore(CI)
	sd := variance.Mul(Z)
	UP := mid.Add(sd)
	LOWER := mid.Sub(sd)

	B := LOW.Gt(LOWER).And(HIGH.Lt(UP))
	df = pandas.NewDataFrame(df.Col("date"))
	ob := pandas.NewSeries(pandas.SERIES_TYPE_BOOL, "cib", B)
	df = df.Join(ob)
	return df
}
