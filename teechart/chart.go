// Package teechart creates charts from Go using the native TeeChart library.
package teechart

// Chart describes a chart. Build it with New() and the fluent setters,
// then call Render() or RenderBytes().
type Chart struct {
	args renderArgs
}

// New returns a Chart with default settings.
func New() *Chart {
	return &Chart{
		args: renderArgs{
			width:   0, // 0 = let the native library decide
			height:  0,
			bufSize: 0,
		},
	}
}

// X sets the X values. If not supplied, x stays empty (the library
// decides what to do), it is not filled with 0..n-1.
func (c *Chart) X(v ...float64) *Chart {
	c.args.x = append([]float64(nil), v...)
	return c
}

// Y sets the Y values.
func (c *Chart) Y(v ...float64) *Chart {
	c.args.y = append([]float64(nil), v...)
	return c
}

// Z sets the Z values (3D series).
func (c *Chart) Z(v ...float64) *Chart {
	c.args.z = append([]float64(nil), v...)
	return c
}

// Texts sets the labels of the points.
func (c *Chart) Texts(v ...string) *Chart {
	c.args.texts = append([]string(nil), v...)
	return c
}

// Style sets the series style, for example "pie", "line" or "bar".
func (c *Chart) Style(s string) *Chart {
	c.args.style = s
	return c
}

// Title sets the chart title.
func (c *Chart) Title(s string) *Chart {
	c.args.title = s
	return c
}

// Footer sets the chart footer.
func (c *Chart) Footer(s string) *Chart {
	c.args.footer = s
	return c
}

// Size sets the output size in pixels.
func (c *Chart) Size(width, height int) *Chart {
	c.args.width = int32(width)
	c.args.height = int32(height)
	return c
}

// Output sets the output type, as understood by the native library.
// Use "modal" or "window" to show the chart in a GUI window, or "png" /
// "jpeg" to render to an in-memory image (retrieve with RenderBytes).
func (c *Chart) Output(s string) *Chart {
	c.args.outputType = s
	return c
}

// Horizontal draws the series horizontally (for example horizontal bars).
func (c *Chart) Horizontal(b bool) *Chart {
	c.args.horiz = b
	return c
}

// Commander shows the chart editor/commander UI.
func (c *Chart) Commander(b bool) *Chart {
	c.args.commander = b
	return c
}

// OpenGL enables OpenGL rendering.
func (c *Chart) OpenGL(b bool) *Chart {
	c.args.opengl = b
	return c
}

// Render draws the chart using the native library.
// When Output is "modal" or "window" this call blocks until the window is closed.
func (c *Chart) Render() error {
	// Work on a copy so the Chart can be reused after Render.
	a := c.args
	return render(&a)
}

// RenderBytes renders the chart to an in-memory image and returns the raw
// bytes. Call Output("png") or Output("jpeg") before RenderBytes so the
// native library produces image data instead of opening a window.
//
//	img, err := teechart.New().
//	    Y(2.5, 3.2, 2.8, 4.1).
//	    Texts("One", "Two", "Three", "Four").
//	    Style("pie").
//	    Title("TeeChart for Go").
//	    Output("png").
//	    Size(800, 500).
//	    RenderBytes()
func (c *Chart) RenderBytes() ([]byte, error) {
	a := c.args
	if err := render(&a); err != nil {
		return nil, err
	}
	// a.bufSize was written back by the DLL with the image byte count.
	return getBytes(a.bufSize)
}
