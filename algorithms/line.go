package algorithms

import "math"

type Line struct {
	m float64 // 斜率
	c float64 // 截距
}

type Point struct {
	x float64
	y float64
}

func calculateDistance(point Point, line Line) float64 {
	A := -line.m
	B := float64(1)
	C := -line.c

	numerator := math.Abs(A*point.x + B*point.y + C)
	denominator := math.Sqrt(math.Pow(A, 2) + math.Pow(B, 2))
	distance := numerator / denominator
	return distance
}

func calculateEquidistantLine(line Line, point Point) Line {
	newC := point.y - line.m*point.x

	newLine := Line{m: line.m, c: newC}
	return newLine
}

func calculateEquidistantLine2(line1 Line, point Point) Line {
	m2 := -1 / line1.m // line2的斜率与line1的斜率相互垂直

	// 计算line1和point之间的垂直距离
	distance := math.Abs(line1.m*point.x-point.y+line1.c) / math.Sqrt(math.Pow(line1.m, 2)+1)

	// 计算line2的截距c2
	c2 := point.y - m2*point.x + distance*math.Sqrt(math.Pow(m2, 2)+1)

	line2 := Line{m: line1.m, c: c2}
	return line2
}

// 已知两个点, 计算直线方程
func calculateLineEquation(point1, point2 Point) Line {
	m := (point2.y - point1.y) / (point2.x - point1.x)
	c := point1.y - m*point1.x

	line := Line{m: m, c: c}
	return line
}
