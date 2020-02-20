package pkg

import "fmt"

// Hello says hello.
func Hello() error {
	fmt.Println("Hello go mod!")
	return nil
}

// Bye says bye.
func Bye() error {
	fmt.Println("Bye go mod!")
	return nil
}
