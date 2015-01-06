package main

import "fmt"

func ExampleValidateId() {
	fmt.Println(validateId("D9L6MbPfE4"))
	fmt.Println(validateId("ABZdez09-_"))
	fmt.Println(validateId("Abc"))
	fmt.Println(validateId("Abc?q=1235"))
	fmt.Println(validateId("../../file"))
	fmt.Println(validateId("Heya世界"))

	// Output:
	// <nil>
	// <nil>
	// id length is 3 instead of 10
	// id contains unexpected character '?'
	// id contains unexpected character '.'
	// id contains unexpected character '\u00e4'
}
