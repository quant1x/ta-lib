package plot

import (
	"fmt"
	"gitee.com/quant1x/num"
	"testing"
)

func TestSeriesIndex(t *testing.T) {
	n := 100
	indexes := num.Range[float64](n)
	fmt.Println(indexes)
}
