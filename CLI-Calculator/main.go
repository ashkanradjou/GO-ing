package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func add(a, b float64) float64 {
	return a + b
}

func subtract(a, b float64) float64 {
	return a - b
}

func multiply(a, b float64) float64 {
	return a * b
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

func processInput(input string) {
	// Split the input into space-separated components
	parts := strings.Fields(input)

	if len(parts) != 3 {
		fmt.Println("Invalid input. Expected format: <num1> <operator> <num2>")
		return
	}

	// Convert the first number to float
	num1, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		fmt.Println("Invalid first number:", parts[0])
		return
	}

	// Convert the second number to float
	num2, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		fmt.Println("Invalid second number:", parts[2])
		return
	}

	// Identify the operator and perform the operation
	switch parts[1] {
	case "+":
		fmt.Printf("Result: %.2f\n", add(num1, num2))
	case "-":
		fmt.Printf("Result: %.2f\n", subtract(num1, num2))
	case "*":
		fmt.Printf("Result: %.2f\n", multiply(num1, num2))
	case "/":
		result, err := divide(num1, num2)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Result: %.2f\n", result)
		}
	default:
		fmt.Println("Invalid operator. Supported operators are +, -, *, /.")
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please provide an operation.")
		return
	}

	input := os.Args[1]

	processInput(input)
}
