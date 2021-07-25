package elsa

import (
	"os"
	"path"
	"path/filepath"

	"shanhu.io/virgo/dock"
)

type env struct {
	dock     *dock.Client
	srcDir   string
	goSrcDir string
	outDir   string

	goVersion      string
	dockerRegistry string
	sshKnownHosts  string
}

func (e *env) prepareOut(ps ...string) (string, error) {
	p := e.out(ps...)
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return p, nil
}

func (e *env) out(ps ...string) string {
	if len(ps) == 0 {
		return e.outDir
	}
	p := path.Join(ps...)
	return filepath.Join(e.outDir, filepath.FromSlash(p))
}

func (e *env) src(ps ...string) string {
	if len(ps) == 0 {
		return e.srcDir
	}
	p := path.Join(ps...)
	return filepath.Join(e.srcDir, filepath.FromSlash(p))
}

func (e *env) gopath() string { return e.src("go") }

func (e *env) goSrc() string {
	if e.goSrcDir != "" {
		return e.goSrcDir
	}
	return filepath.Join(e.gopath(), "src")
}

func (e *env) docker() *dock.Client { return e.dock }

func (e *env) dockerName(s string) string {
	if e.dockerRegistry == "" {
		return s
	}
	return path.Join(e.dockerRegistry, s)
}
