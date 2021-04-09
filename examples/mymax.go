package main

import (
	"fmt"
)





func main() {
	intVal := myMax_gonerics_c1odp2fcsae3rk5958n0(10, 20)
	floatVal := myMax_gonerics_c1odp2fcsae3rk5958ng(3.14, 5.64)

	Printer_gonerics_c1odp2fcsae3rk5958o0(25, "Hello World")

	fmt.Printf("%+v\n", intVal)
	fmt.Printf("%+v\n", floatVal)
}
func myMax_gonerics_c1odp2fcsae3rk5958n0(x int,y int) (int) {
							if x > y {
		return x
	}

	return y
					}
func myMax_gonerics_c1odp2fcsae3rk5958ng(x float64,y float64) (float64) {
							if x > y {
		return x
	}

	return y
					}
func Printer_gonerics_c1odp2fcsae3rk5958o0(size int,data string) (int) {
							fmt.Printf("%+v\n", data)

	return 1
					}
