// Copyright (C) 2021  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package elsa

import (
	"os"
	"path"
	"path/filepath"
	"strings"

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
		if e.goSrcDir == "@" {
			if sys := systemGoSrc(); sys != "" {
				return sys
			}
		} else {
			return e.goSrcDir
		}
	}
	return filepath.Join(e.gopath(), "src")
}

func (e *env) docker() *dock.Client { return e.dock }

func (e *env) dockerName(s string) string {
	if e.dockerRegistry == "" {
		return s
	}
	const prefix = "docker/"
	if strings.HasPrefix(s, prefix) {
		s = strings.TrimPrefix(s, prefix)
	}
	return path.Join(e.dockerRegistry, s)
}
