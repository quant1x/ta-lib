package algorithms

import "math"

// Polyfit 用最小二乘法拟合多项式系数
//
//	x: 自变量数组，y: 因变量数组，degree: 多项式阶数
//	返回系数数组（从高阶到低阶），例如 degree=1 时返回 [斜率, 截距]
func Polyfit(x, y []float64, degree int) []float64 {
	// 构造设计矩阵
	n := len(x)
	m := degree + 1
	X := make([][]float64, n)
	for i := 0; i < n; i++ {
		X[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			X[i][j] = math.Pow(x[i], float64(j))
		}
	}

	// 计算 X^T * X 和 X^T * Y
	XTX := make([][]float64, m)
	for i := 0; i < m; i++ {
		XTX[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			for k := 0; k < n; k++ {
				XTX[i][j] += X[k][i] * X[k][j]
			}
		}
	}

	XTY := make([]float64, m)
	for i := 0; i < m; i++ {
		for k := 0; k < n; k++ {
			XTY[i] += X[k][i] * y[k]
		}
	}

	// 解线性方程组 XTX * Coef = XTY
	return gauss(XTX, XTY)
}

// 高斯消元法解线性方程组
func gauss(A [][]float64, b []float64) []float64 {
	n := len(A)
	for i := 0; i < n; i++ {
		// 寻找主元
		maxRow := i
		for j := i; j < n; j++ {
			if math.Abs(A[j][i]) > math.Abs(A[maxRow][i]) {
				maxRow = j
			}
		}
		A[i], A[maxRow] = A[maxRow], A[i]
		b[i], b[maxRow] = b[maxRow], b[i]

		// 消元
		for j := i + 1; j < n; j++ {
			c := A[j][i] / A[i][i]
			for k := i; k < n; k++ {
				A[j][k] -= c * A[i][k]
			}
			b[j] -= c * b[i]
		}
	}

	// 回代
	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		x[i] = b[i]
		for j := i + 1; j < n; j++ {
			x[i] -= A[i][j] * x[j]
		}
		x[i] /= A[i][i]
	}
	return x
}
