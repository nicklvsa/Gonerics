package main

import (
	"fmt"
)





func main() {
	largestInt := myMax_gonerics_c3b6goncsae4676p2q7g(20, 60)
	fmt.Printf("%+v\n", largestInt)

	largestFloat := myMax_gonerics_c3b6goncsae4676p2q80(3.14, 8.25)
	fmt.Printf("%+v\n", largestFloat)
}
func myMax_gonerics_c3b6goncsae4676p2q7g(num0 int,num1 int) (int) {
							if num0 > num1 {
		return num0
	}

	return num1
					}
func myMax_gonerics_c3b6goncsae4676p2q80(num0 float64,num1 float64) (float64) {
							if num0 > num1 {
		return num0
	}

	return num1
					}
