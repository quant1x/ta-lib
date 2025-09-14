package plot

import (
	"gitee.com/quant1x/num"
	"gitee.com/quant1x/pandas"
)

func SeriesIndex(s pandas.Series) []float64 {
	n := s.Len()
	//indexes := make([]float64, n)
	//for i:= 0; i < n ; i++ {
	//	stat.Range()
	//}
	indexes := num.Range[float64](n)
	return indexes
}
