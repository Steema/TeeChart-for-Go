# TeeChart-for-Go

TeeChart charting / plotting / graphing library for Go, powered by the native
`teechart.dll` (Windows) or `libteechart.so` / `libteechart.dylib`
(Linux / macOS).

---

<img width="1810" height="955" alt="image" src="https://github.com/user-attachments/assets/5bb0ee2c-fc3b-4394-8bb9-a604649b0cfe" />


### Notes

Tested on Windows with Go 1.27.

The 64-bit DLL (`libs/win64/teechart.dll`) is included in the repository.
On Linux and macOS you must supply the shared library separately and point
`TEECHART_LIB` to its location.

---

### Quick start

```go
go run ./examples/simple
```

This creates a pie chart in a modal window:

```go
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
```

---

### Rendering to an image file

Use `Output("png")` (or `"jpeg"`) together with `RenderBytes()` to get the
chart as a byte slice and write it to disk:

```go
img, err := teechart.New().
    Y(2.5, 3.2, 2.8, 4.1).
    Style("pie").
    Title("TeeChart for Go").
    Output("png").
    Size(800, 500).
    RenderBytes()
if err != nil {
    log.Fatal(err)
}
os.WriteFile("chart.png", img, 0644)
```

Run the bitmap example:

```
go run ./examples/bitmap
```

---

### Library search order

The package locates the native library automatically:

1. `TEECHART_LIB` environment variable (file or folder)
2. Next to the executable (and in `libs/<platform>/` below it)
3. `libs/<platform>/` relative to the current working directory
4. `libs/<platform>/` inside the module source folder

Set `TEECHART_LIB` to override all of the above:

```
set TEECHART_LIB=C:\path\to\teechart.dll
```

---

### Fluent API reference

| Method | Description |
|--------|-------------|
| `New()` | Create a new chart with default settings |
| `.X(v ...float64)` | Set X values |
| `.Y(v ...float64)` | Set Y values |
| `.Z(v ...float64)` | Set Z values (3-D series) |
| `.Texts(v ...string)` | Set point labels |
| `.Style(s string)` | Series style: `"pie"`, `"bar"`, `"line"`, … |
| `.Title(s string)` | Chart title |
| `.Footer(s string)` | Chart footer |
| `.Size(w, h int)` | Output size in pixels |
| `.Output(s string)` | Output type: `"modal"`, `"window"`, `"png"`, `"jpeg"` |
| `.Horizontal(b bool)` | Draw series horizontally |
| `.Commander(b bool)` | Show the chart editor UI |
| `.OpenGL(b bool)` | Enable OpenGL rendering |
| `.Render() error` | Draw the chart (window or image file inside DLL) |
| `.RenderBytes() ([]byte, error)` | Render and return raw image bytes |
| `teechart.GetBytes()` | Retrieve bytes from the last Render call |
