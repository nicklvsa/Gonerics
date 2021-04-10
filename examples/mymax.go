package main

import (
	"fmt"
)





func main() {
	intVal := myMax_gonerics_c1ogoifcsae3dp7gs3o0(10, 20)
	floatVal := myMax_gonerics_c1ogoifcsae3dp7gs3og(3.14, 5.64)

	Printer_gonerics_c1ogoifcsae3dp7gs3p0("Hello World")

	fmt.Printf("%+v\n", intVal)
	fmt.Printf("%+v\n", floatVal)
}
func myMax_gonerics_c1ogoifcsae3dp7gs3o0(x int,y int) (int) {
							if x > y {
		return x
	}

	return y
					}
func myMax_gonerics_c1ogoifcsae3dp7gs3og(x float64,y float64) (float64) {
							if x > y {
		return x
	}

	return y
					}
func Printer_gonerics_c1ogoifcsae3dp7gs3p0(data string)  {
							fmt.Printf("%+v\n", data)
					}
