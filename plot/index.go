package plot

import "gitee.com/quant1x/pandas/stat"

func SeriesIndex(s stat.Series) []float64 {
	n := s.Len()
	//indexes := make([]float64, n)
	//for i:= 0; i < n ; i++ {
	//	stat.Range()
	//}
	indexes := stat.Range[float64](n)
	return indexes
}
