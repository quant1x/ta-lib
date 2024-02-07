package plot

import (
	"fmt"
	"gitee.com/quant1x/pandas/stat"
	"testing"
)

func TestSeriesIndex(t *testing.T) {
	n := 100
	indexes := stat.Range[float64](n)
	fmt.Println(indexes)
}
