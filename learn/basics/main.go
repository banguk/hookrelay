package main

import "fmt"

// Package-level variable: must use var, := is not allowed here.
var appName = "hookrelay"

// Untyped constant: Go picks the type where it is used.
const maxRetries = 5

func add(a int, b int) int {
	return a + b
}

// Parameters of the same type can share one type annotation.
func divmod(a, b int) (int, int) {
	return a / b, a % b
}

// Named return values: declared up front, bare return sends them back.
func minmax(a, b int) (min, max int) {
	if a < b {
		min, max = a, b
	} else {
		min, max = b, a
	}
	return
}

func main() {
	var greeting string = "Hello from"

	version := 0.1

	var attempts int
	var delivered bool
	var lastError string

	fmt.Println(greeting, appName, version)
	fmt.Println("zero values:", attempts, delivered, lastError == "")

	sum := add(2, 3)
	fmt.Println("add:", sum)

	q, r := divmod(17, 5)
	fmt.Println("divmod:", q, r)

	_, remainder := divmod(maxRetries, 2)
	fmt.Println("remainder only:", remainder)

	lo, hi := minmax(9, 4)
	fmt.Println("minmax:", lo, hi)

	fmt.Printf("attempts=%d appName=%s version=%v type=%T\n", attempts, appName, version, version)

}
