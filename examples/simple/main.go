package main

import (
	"fmt"
	"os"

	"github.com/Steema/TeeChart-for-Go/teechart"
)

func main() {
	err := teechart.New().
		Y(2.5, 3.2, 2.8, 4.1).
		Texts("One", "Two", "Three", "Four").
		Style("pie").
		Title("TeeChart for Go").
		Render()

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}