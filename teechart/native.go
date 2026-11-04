package teechart

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

// renderArgs holds everything TeeChartRender needs, in Go-friendly types.
type renderArgs struct {
	x, y, z    []float64
	texts      []string
	outputType string
	width      int32
	height     int32
	style      string
	title      string
	footer     string
	horiz      bool
	commander  bool
	opengl     bool
	bufSize    int32 // written back by the DLL with the image byte count
}

var (
	loadOnce sync.Once
	loadErr  error

	// Mirrors the Delphi declaration exactly: every parameter is a pointer.
	fnRender func(
		nx *int32, x *float64,
		ny *int32, y *float64,
		nz *int32, z *float64,
		ntexts *int32, texts **byte,
		outputType **byte,
		width, height *int32,
		style, title, footer **byte,
		isHoriz, commander, opengl *int32,
		bufSize *int32,
	) int32

	// fnGetBytes copies the last rendered image into a caller-supplied buffer.
	// Delphi: function TeeChartGetBytes(Buf: Pointer; var Size: Integer): Boolean
	// Buf must be pre-allocated to at least Size bytes; on return Size holds the
	// actual number of bytes written.
	fnGetBytes func(buf unsafe.Pointer, size *int32) bool
)

func load() error {
	loadOnce.Do(func() {
		path, err := findLibrary()
		if err != nil {
			loadErr = err
			return
		}
		handle, err := openLibrary(path)
		if err != nil {
			loadErr = fmt.Errorf("teechart: cannot load %q: %w", path, err)
			return
		}
		defer func() {
			if r := recover(); r != nil {
				loadErr = fmt.Errorf("teechart: cannot bind TeeChartRender in %q: %v", path, r)
			}
		}()
		purego.RegisterLibFunc(&fnRender, handle, "TeeChartRender")
		purego.RegisterLibFunc(&fnGetBytes, handle, "TeeChartGetBytes")
	})
	return loadErr
}

func cstring(s string) []byte {
	b := make([]byte, len(s)+1)
	copy(b, s)
	return b
}

func b2i(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

// dptr never returns nil: an empty slice gives a pointer to a dummy value
// (count 0), so the Delphi side can't crash dereferencing nil.
func dptr(v []float64) *float64 {
	if len(v) == 0 {
		return new(float64)
	}
	return &v[0]
}

func render(a *renderArgs) error {
	if err := load(); err != nil {
		return err
	}

	nx, ny, nz := int32(len(a.x)), int32(len(a.y)), int32(len(a.z))
	nt := int32(len(a.texts))

	// texts: array of char*
	tbufs := make([][]byte, len(a.texts))
	tptrs := make([]*byte, len(a.texts)+1) // +1 keeps it non-nil when empty
	for i, t := range a.texts {
		tbufs[i] = cstring(t)
		tptrs[i] = &tbufs[i][0]
	}

	// single strings: PPAnsiChar = pointer to one char*
	ot, st, ti, fo := cstring(a.outputType), cstring(a.style), cstring(a.title), cstring(a.footer)
	otp, stp, tip, fop := &ot[0], &st[0], &ti[0], &fo[0]

	horiz, cmd, gl := b2i(a.horiz), b2i(a.commander), b2i(a.opengl)
	w, h, buf := a.width, a.height, a.bufSize

	rc := fnRender(
		&nx, dptr(a.x),
		&ny, dptr(a.y),
		&nz, dptr(a.z),
		&nt, &tptrs[0],
		&otp,
		&w, &h,
		&stp, &tip, &fop,
		&horiz, &cmd, &gl,
		&buf,
	)

	runtime.KeepAlive(a)
	runtime.KeepAlive(tbufs)
	runtime.KeepAlive(tptrs)
	runtime.KeepAlive(ot)
	runtime.KeepAlive(st)
	runtime.KeepAlive(ti)
	runtime.KeepAlive(fo)

	a.bufSize = buf // DLL writes back the actual image byte count

	if rc != 0 {
		return fmt.Errorf("teechart: render failed (code %d)", rc)
	}
	return nil
}

// getBytes copies the image bytes from the DLL's internal buffer into a new
// Go slice.  byteCount must be the value written back into bufSize by render().
// Returns nil when byteCount is 0 (e.g. the output was a window, not an image).
func getBytes(byteCount int32) ([]byte, error) {
	if err := load(); err != nil {
		return nil, err
	}
	if byteCount <= 0 {
		return nil, nil
	}

	out := make([]byte, byteCount)
	size := byteCount
	ok := fnGetBytes(unsafe.Pointer(&out[0]), &size)
	if !ok {
		return nil, fmt.Errorf("teechart: TeeChartGetBytes failed")
	}
	return out[:size], nil
}
