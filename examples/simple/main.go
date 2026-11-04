// simple shows the most basic TeeChart-for-Go usage: render a pie chart
// in a modal window (the window closes when you click OK or close it).
//
// Usage:
//
//	go run ./examples/simple
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
		Output("modal").
		Size(800, 500).
		Render()

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Println("TeeChartRender correctly executed")
}
