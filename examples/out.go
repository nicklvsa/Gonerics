package main

import (
	"fmt"
)

type MyStruct struct {
	Name string
	Age int
}







func main() {
	test()
	cool_gonerics_c1mehfncsae4pv12p6ng(15, 19)

	my := MyStruct{
		Name: "Bob",
		Age: 30,
	}

	something_gonerics_c1mehfncsae4pv12p6o0("hello world", &my)
}

func test() {
	example_gonerics_c1mehfncsae4pv12p6og("yo", nil, nil)
	example_gonerics_c1mehfncsae4pv12p6og("hello", nil, nil)
}
func cool_gonerics_c1mehfncsae4pv12p6ng(num int,age int) (string,int) {
							return "hello world", age
					}
func something_gonerics_c1mehfncsae4pv12p6o0(arg0 string,arg1 *MyStruct)  {
							fmt.Printf("%+v\n", arg1.Name)
					}
func example_gonerics_c1mehfncsae4pv12p6og(input string,another error,yo *string) (a *string,b string) {
							fmt.Println(input)
	fmt.Println(another)

	return yo, input
					}
func example_gonerics_c1mehfncsae4pv12p6p0(input string,another error,yo *string) (a *string,b string) {
							fmt.Println(input)
	fmt.Println(another)

	return yo, input
					}
