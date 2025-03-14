package linear

import (
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/pandas"
)

// CurveRegression 曲线回归
//
//	curve regression
//	https://blog.csdn.net/llittleSun/article/details/115045660
//	默认只预1个周期
//	argv 第一个参数为数据长度, 即周期数, 默认为S的长度
func CurveRegression(S pandas.Series, argv ...int) pandas.Series {
	N := S.Len()
	if len(argv) > 0 {
		N = argv[0]
	}
	if N > S.Len() {
		N = S.Len()
	}

	y := S.Select(api.RangeFinite(-N)).DTypes()
	x := num.Arange[num.DType](1, float64(N)+1, 1)
	t1 := num.Pow[num.DType](x, 2)
	t2 := x
	t3 := num.Ones[num.DType](x)

	A := num.Concat1D(t1, t2, t3)
	T := num.Transpose2D(A)

	w0 := num.Dot2D[num.DType](T, A)
	w1 := num.Inverse(w0)

	w2 := num.Dot2D(w1, T)
	W := num.Dot2D1[num.DType](w2, y)

	d1 := num.Arange[num.DType](1, num.DType(N)+2, 1)

	d21 := num.Pow(d1, 2)
	d2 := pandas.ToSeries(d21...).Mul(W[0])
	d3 := pandas.ToSeries(d1...).Mul(W[1]).Add(W[2])

	D := d2.Add(d3)
	return D
}
