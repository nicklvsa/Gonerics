package main

import (
	"fmt"
)







func main() {
	intVal := myMax_gonerics_c1p2c97bm9vsor1rrpg0(10, 20)
	floatVal := myMax_gonerics_c1p2c97bm9vsor1rrpgg(3.14, 5.64)

	Printer_gonerics_c1p2c97bm9vsor1rrph0("Hello World")

	fmt.Printf("%+v\n", intVal)
	fmt.Printf("%+v\n", floatVal)
}
func myMax_gonerics_c1p2c97bm9vsor1rrpg0(x int,y int) (int) {
							if x > y {
		return x
	}

	return y
					}
func myMax_gonerics_c1p2c97bm9vsor1rrpgg(x float64,y float64) (float64) {
							if x > y {
		return x
	}

	return y
					}
func Printer_gonerics_c1p2c97bm9vsor1rrph0(data string)  {
							fmt.Printf("%+v\n", data)
					}
