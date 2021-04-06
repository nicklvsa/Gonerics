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
	cool_gonerics_c1lseh7bm9vl11ce6q5g(15, 19)

	my := MyStruct{
		Name: "Bob",
		Age: 30,
	}

	something_gonerics_c1lseh7bm9vl11ce6q60("hello world", &my)
}

func test() {
	example_gonerics_c1lseh7bm9vl11ce6q6g("yo", nil, nil)
	example_gonerics_c1lseh7bm9vl11ce6q70("hello", nil, nil)
}
func cool_gonerics_c1lseh7bm9vl11ce6q5g(num int,age int) (string,int) {
							return "hello world", age
					}
func something_gonerics_c1lseh7bm9vl11ce6q60(arg0 string,arg1 *MyStruct)  {
							fmt.Printf("%+v\n", arg1.Name)
					}
func example_gonerics_c1lseh7bm9vl11ce6q6g(input string,another error,yo *string) (a *string,b string) {
							fmt.Println(input)
	fmt.Println(another)

	return yo, input
					}
func example_gonerics_c1lseh7bm9vl11ce6q70(input string,another error,yo *string) (a *string,b string) {
							fmt.Println(input)
	fmt.Println(another)

	return yo, input
					}
