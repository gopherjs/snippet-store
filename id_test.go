package main

import "fmt"

func ExampleValidateID() {
	fmt.Println(validateID("D9L6MbPfE4"))
	fmt.Println(validateID("ABZdez09-_"))
	fmt.Println(validateID("Abc"))
	fmt.Println(validateID("Abc?q=1235"))
	fmt.Println(validateID("../../file"))
	fmt.Println(validateID("Heya世界"))

	// Output:
	// <nil>
	// <nil>
	// id length is 3 instead of 10
	// id contains unexpected character '?'
	// id contains unexpected character '.'
	// id contains unexpected character '\u00e4'
}
