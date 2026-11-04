// bitmap demonstrates rendering a TeeChart chart to a PNG byte slice
// using the TeeChartGetBytes DLL export, then saving it to disk.
//
// Usage:
//
//	go run ./examples/bitmap
package main

import (
	"fmt"
	"os"

	"github.com/Steema/TeeChart-for-Go/teechart"
)

func main() {
	img, err := teechart.New().
		Y(2.5, 3.2, 2.8, 4.1).
		Texts("One", "Two", "Three", "Four").
		Style("pie").
		Title("TeeChart for Go").
		Size(800, 500).
		Output("png").
		RenderBytes()

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if len(img) == 0 {
		fmt.Fprintln(os.Stderr, "no image bytes returned by the DLL")
		os.Exit(1)
	}

	const outFile = "chart.png"
	if err := os.WriteFile(outFile, img, 0644); err != nil {
		fmt.Fprintln(os.Stderr, "error writing file:", err)
		os.Exit(1)
	}

	fmt.Printf("Chart saved to %s (%d bytes)\n", outFile, len(img))
}

