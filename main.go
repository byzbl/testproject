package main

import (
	"fmt"
	"reflect"
)

func main() {

	type S struct {
		F string `species:"gopher" color:"blue"`
	}

	s := S{}
	st := reflect.TypeOf(s)
	field, _ := st.FieldByName("F")
	//field := st.Field(0)
	fmt.Println(field.Tag.Get("color"), field.Tag.Get("species"))
	//fmt.Println("111hello git!")

	fmt.Println("test")
	fmt.Println("test2")
	fmt.Println("test2")

}
