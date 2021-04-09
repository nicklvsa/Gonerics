package main

import (
	"fmt"
)





func main() {
	intVal := myMax_gonerics_c1odqcvcsae2d60tbigg(10, 20)
	floatVal := myMax_gonerics_c1odqcvcsae2d60tbih0(3.14, 5.64)

	Printer_gonerics_c1odqcvcsae2d60tbihg(25, "Hello World")

	fmt.Printf("%+v\n", intVal)
	fmt.Printf("%+v\n", floatVal)
}
func myMax_gonerics_c1odqcvcsae2d60tbigg(x int,y int) (int) {
							if x > y {
		return x
	}

	return y
					}
func myMax_gonerics_c1odqcvcsae2d60tbih0(x float64,y float64) (float64) {
							if x > y {
		return x
	}

	return y
					}
func Printer_gonerics_c1odqcvcsae2d60tbihg(size int,data string) (int) {
							fmt.Printf("%+v\n", data)

	return 1
					}
