package main

import (
	"fmt"
)

func main() {
	intVal := myMax_gonerics_c1ns7o7bm9vlv77rje6g(10, 20)
	floatVal := myMax_gonerics_c1ns7o7bm9vlv77rje70(3.14, 5.64)

	fmt.Printf("%+v\n", intVal)
	fmt.Printf("%+v\n", floatVal)
}
func myMax_gonerics_c1ns7o7bm9vlv77rje6g(x int, y int) int {
	if x > y {
		return x
	}

	return y
}
func myMax_gonerics_c1ns7o7bm9vlv77rje70(x float64, y float64) float64 {
	if x > y {
		return x
	}

	return y
}
