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
	"log"
	"os"

	"shanhu.io/misc/errcode"
	"shanhu.io/virgo/dock"
)

// Config provide the configuration to start a builder.
type Config struct {
	Src   string // Source directory
	GoSrc string // Alternative golang source directory
	Out   string // Output directory

	DockerRegistry string // Docker registry for output tagging.
	GoVersion      string // File to read golang's version.
	SSHKnownHosts  string // File of the SSH known hosts list.
}

// Builder builds stuff.
type Builder struct {
	env *env
}

// NewBuilder creates a new builder that builds stuff.
func NewBuilder(config *Config) *Builder {
	env := &env{
		dock:           dock.NewUnixClient(""),
		srcDir:         config.Src,
		goSrcDir:       config.GoSrc,
		outDir:         config.Out,
		goVersion:      config.GoVersion,
		dockerRegistry: config.DockerRegistry,
		sshKnownHosts:  config.SSHKnownHosts,
	}

	return &Builder{
		env: env,
	}
}

func (b *Builder) buildBase(dockers []*baseDocker) error {
	if err := os.MkdirAll(b.env.out(), 0700); err != nil {
		return errcode.Annotate(err, "make output dir")
	}

	for _, d := range dockers {
		log.Printf("build %s", d.name)
		if err := d.build(b.env, d.name); err != nil {
			return errcode.Annotatef(err, "build %s", d.name)
		}
	}
	return nil
}

// BuildBaseNodeJS builds the nodejs base docker.
func (b *Builder) BuildBaseNodeJS() error {
	dockers := []*baseDocker{dockerNodejs}
	return b.buildBase(dockers)
}

// BuildBase builds all builtin base dockers.
func (b *Builder) BuildBase() error {
	dockers := []*baseDocker{
		dockerAlpine,
		dockerHatch,
		dockerGolang,
		dockerForge,
		dockerNodejs,
		dockerApp,
	}

	if err := b.buildBase(dockers); err != nil {
		return err
	}

	log.Println("all base docker built")
	return nil
}

// BuildBin builds alpine binaries.
func (b *Builder) BuildBin(pkgs []string) error {
	alpine := newAlpine(b.env, b.env.dockerName(dockerForge.name))
	return alpine.build(pkgs)
}

// BuildModBin builds golang alpine binaries in module-aware mode.
func (b *Builder) BuildModBin(dir string, pkgs []string) error {
	alpine := newAlpine(b.env, b.env.dockerName(dockerForge.name))
	return alpine.buildMod(dir, pkgs)
}

// BuildNodeJS builds an nodejs npm package.
func (b *Builder) BuildNodeJS(dir string, js *NodeJS) error {
	nodeJS := newNodeJS(b.env, b.env.dockerName(dockerNodejs.name))
	return nodeJS.build(dir, js)
}

// BuildDocker builds a docker.
func (b *Builder) BuildDocker(dir string, saveName bool) error {
	d := newDocker(b.env)
	return d.build(dir, saveName)
}

// PullDockers pulls dockers and save them in output.
func (b *Builder) PullDockers(
	dir string, p *DockerPull, opt *DockerPullOptions,
) error {
	return pullDockers(b.env, dir, p, opt)
}

// Src returns the filesystem path to a source file.
func (b *Builder) Src(f string) string { return b.env.src(f) }

// Out returns the filesystme path to an output file.
func (b *Builder) Out(f string) string { return b.env.out(f) }

// GOPATH returns the GOPATH for go language binary building (in non-module
// mode).
func (b *Builder) GOPATH() string { return b.env.gopath() }

// SyncRepos synchronizes the repositories. When sums is nil, it pulls
// from the latest HEAD.
func (b *Builder) SyncRepos(build *Build, sums *BuildSums) (
	*BuildSums, error,
) {
	return syncRepos(b.env, build, sums)
}
