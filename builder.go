// Copyright (C) 2022  Shanhu Tech Inc.
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

package caco3

import (
	"fmt"
	"log"
	"os"

	"shanhu.io/misc/errcode"
	"shanhu.io/text/lexing"
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
func NewBuilder(workDir string, config *Config) *Builder {
	env := &env{
		dock:           dock.NewUnixClient(""),
		workDir:        workDir,
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
func (b *Builder) SyncRepos(ws *Workspace, sums *RepoSums) (
	*RepoSums, error,
) {
	return syncRepos(b.env, ws, sums)
}

// Build builds the given rules.
func (b *Builder) Build(rules []string) []*lexing.Error {
	nodes, nodeMap, errs := loadNodes(b.env, rules)
	if errs != nil {
		return errs
	}
	ctx := &buildContext{nodes: nodeMap, built: make(map[string]string)}
	return b.buildNodes(ctx, nodes)
}

func (b *Builder) buildNodes(
	ctx *buildContext, nodes []*buildNode,
) []*lexing.Error {
	b.env.nodeType = ctx.nodeType
	b.env.ruleType = ctx.ruleType

	for _, n := range nodes {
		if n.typ == nodeSrc {
			log.Printf("%s is a source file", n.name)
			continue
		}
		if _, err := b.buildNode(ctx, n); err != nil {
			return lexing.SingleErr(err)
		}
	}
	return nil
}

func (b *Builder) buildNode(ctx *buildContext, n *buildNode) (
	string, error,
) {
	if digest, ok := ctx.built[n.name]; ok {
		return digest, nil
	}
	digest := ""
	defer func() { ctx.built[n.name] = digest }()

	deps := make(map[string]string)
	for _, dep := range n.deps {
		depNode := ctx.nodes[dep]
		if depNode == nil {
			return "", fmt.Errorf(
				"dep %q for %q not found", dep, n.name,
			)
		}
		d, err := b.buildNode(ctx, depNode)
		if err != nil {
			return "", err
		}
		if d == "" {
			// If any dep is always rebuilding, then this node
			// is also always rebuilding.
			deps = nil
		} else if deps != nil {
			deps[dep] = d
		}
	}

	if deps != nil { // Not always rebuilding, so calculate the digest
		switch n.typ {
		case nodeRule:
			action := &buildAction{
				Deps:     deps,
				RuleType: n.ruleType,
			}
			if meta := n.ruleMeta; meta != nil {
				if meta.digest == "" {
					break
				}
				action.Rule = meta.digest
			}
			d, err := makeDigest("build_action", "", action)
			if err != nil {
				return "", errcode.Annotate(err, "digest build action")
			}
			digest = d
		case nodeSrc:
			stat, err := newSrcFileStat(b.env, n.name)
			if err != nil {
				return "", errcode.Annotatef(err, "stat file %q", n.name)
			}
			d, err := makeDigest("src", "", stat)
			if err != nil {
				return "", errcode.Annotate(err, "digest source file")
			}
			digest = d
		case nodeOut:
			// TODO(h8liu): load output origin digest
		}
	}

	if digest != "" {
		// TODO(h8liu): check build cache here
	}

	if n.rule != nil {
		if n.typ == nodeRule {
			log.Printf("BUILD %s", n.name)
			// TODO(h8liu): better opts
			opts := &buildOpts{
				log: os.Stderr,
				docker: &dockerOpts{
					useBuildCache: true,
				},
			}
			if err := n.rule.build(b.env, opts); err != nil {
				return "", errcode.Annotatef(err, "build %s", n.name)
			}
		}
	}
	return digest, nil
}
