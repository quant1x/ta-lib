package algorithms

import (
	"fmt"
	"testing"
)

func Test_calculateEquidistantLine_basic(t *testing.T) {
	line := Line{m: 2, c: 1}
	point := Point{x: 3, y: 4}

	newLine := calculateEquidistantLine(line, point)
	fmt.Printf("New line equation: y = %.2fx + %.2f\n", newLine.m, newLine.c)
}
