package plot

import (
	"fmt"
	"testing"

	"gitee.com/quant1x/num"
)

func TestSeriesIndex(t *testing.T) {
	n := 100
	indexes := num.Range[float64](n)
	fmt.Println(indexes)
}
