package main

import (
	"fmt"
)

func main() {
	var a = []int{1, 2, 3}
	fmt.Print(a...)
	// fmt.Printf("Decomposed: %s", a...)
}
