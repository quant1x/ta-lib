package plot

import (
	"fmt"
	"testing"

	"github.com/quant1x/num"
)

func TestSeriesIndex(t *testing.T) {
	n := 100
	indexes := num.Range[float64](n)
	fmt.Println(indexes)
}
