package elsa

import (
	"os"
	"path/filepath"
)

func systemGoSrc() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		home := os.Getenv("HOME")
		if home == "" {
			return ""
		}
		return filepath.Join(home, "go", "src")
	}
	return filepath.Join(gopath, "src")
}
