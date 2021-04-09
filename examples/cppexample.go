package main

import (
	"fmt"
)

type Struct0 struct {
	Name string
}

type Struct1 struct {
	Name string
	another int
}





func main() {
	intVal := myMax_gonerics_c1oc02fcsae2qa0vpmsg(10, 20)
	floatVal := myMax_gonerics_c1oc02fcsae2qa0vpmt0(3.14, 5.64)

	st1 := Struct1{
		Name: "The number is:",
	}

	updated := appendToStruct_gonerics_c1oc02fcsae2qa0vpmtg(20, &st1)

	fmt.Printf("%+v\n", intVal)
	fmt.Printf("%+v\n", floatVal)
	fmt.Printf("%+v\n", *updated)
}
func myMax_gonerics_c1oc02fcsae2qa0vpmsg(x int,y int) (int) {
							if x > y {
		return x
	}

	return y
					}
func myMax_gonerics_c1oc02fcsae2qa0vpmt0(x float64,y float64) (float64) {
							if x > y {
		return x
	}

	return y
					}
func appendToStruct_gonerics_c1oc02fcsae2qa0vpmtg(number int64,st *Struct1) (*Struct1) {
							st.Name = fmt.Sprintf("%s %+v", st.Name, number)
	return st
					}
