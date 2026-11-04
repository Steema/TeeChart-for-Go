package teechart

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const envLib = "TEECHART_LIB"

// platformInfo returns the subfolder inside libs/ and the library file name
// for the current OS and architecture.
func platformInfo() (dir, file string, err error) {
	switch runtime.GOOS {
	case "windows":
		if runtime.GOARCH == "amd64" {
			return "win64", "teechart.dll", nil
		}
	case "linux":
		if runtime.GOARCH == "amd64" {
			return "linux64", "libteechart.so", nil
		}
	case "darwin":
		return "macos", "libteechart.dylib", nil
	}
	return "", "", fmt.Errorf("teechart: unsupported platform %s/%s", runtime.GOOS, runtime.GOARCH)
}

// findLibrary locates the native TeeChart library.
//
// Search order:
//  1. TEECHART_LIB environment variable (file or folder)
//  2. next to the executable (and in libs/<platform>/ below it)
//  3. libs/<platform>/ relative to the current working directory
//  4. libs/<platform>/ inside this module's source folder
func findLibrary() (string, error) {
	dir, file, err := platformInfo()
	if err != nil {
		return "", err
	}

	// 1. Explicit path: if set, it must be valid (no silent fallback).
	if p := os.Getenv(envLib); p != "" {
		if st, err := os.Stat(p); err == nil {
			if st.IsDir() {
				p = filepath.Join(p, file)
				if _, err := os.Stat(p); err != nil {
					return "", fmt.Errorf("teechart: %s=%q does not contain %s", envLib, filepath.Dir(p), file)
				}
			}
			return absPath(p), nil
		}
		return "", fmt.Errorf("teechart: %s=%q not found", envLib, p)
	}

	var candidates []string

	// 2. Next to the executable.
	if exe, err := os.Executable(); err == nil {
		if real, err := filepath.EvalSymlinks(exe); err == nil {
			exe = real
		}
		exeDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(exeDir, file),
			filepath.Join(exeDir, "libs", dir, file),
		)
	}

	// 3. Current working directory (makes "go run ./examples/simple" work).
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(wd, "libs", dir, file))
	}

	// 4. Module source folder (works when the module is in the Go module
	// cache after "go get"; not available if built with -trimpath).
	if _, thisFile, _, ok := runtime.Caller(0); ok {
		root := filepath.Dir(filepath.Dir(thisFile)) // teechart/ -> repo root
		candidates = append(candidates, filepath.Join(root, "libs", dir, file))
	}

	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return absPath(c), nil
		}
	}

	return "", fmt.Errorf("teechart: native library %q not found; set %s or place it in libs/%s/\ntried:\n  %s",
		file, envLib, dir, strings.Join(candidates, "\n  "))
}

func absPath(p string) string {
	if a, err := filepath.Abs(p); err == nil {
		return a
	}
	return p
}