package labs

import (
	"fmt"
	"testing"
)

func TestVPVR(t *testing.T) {
	// 测试数据
	close := []float64{10.0, 10.5, 11.0, 10.8, 10.2, 10.7, 10.9}
	volume := []float64{1000, 2000, 1500, 3000, 1200, 2500, 1800}

	// 计算VPVR
	result, err := VPVR(close, volume, 5)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Max Volume Price Bin: %.2f\n", result) // 输出: 10.80
}
