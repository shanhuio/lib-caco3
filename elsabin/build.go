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

package elsabin

import (
	"log"
	"path"

	"shanhu.io/elsa"
	"shanhu.io/misc/errcode"
)

type targets struct {
	m     map[string]*elsa.BuildStep
	steps []*elsa.BuildStep
}

func newTargets(b *elsa.Build) *targets {
	m := make(map[string]*elsa.BuildStep)
	for _, step := range b.Steps {
		m[step.Name] = step
	}
	return &targets{
		m:     m,
		steps: b.Steps,
	}
}

func buildDockers(b *elsa.Builder, dir string, docks []string) error {
	for _, d := range docks {
		if err := b.BuildDocker(path.Join(dir, d)); err != nil {
			return errcode.Annotatef(err, "build docker %q", d)
		}
	}
	return nil
}

type buildOptions struct {
	DockerPull *elsa.DockerPullOptions
}

func buildStep(
	b *elsa.Builder, step *elsa.BuildStep, opts *buildOptions,
) error {
	log.Printf("build %s", step.Name)
	dir := step.Dir
	if dir == "" {
		dir = step.Name
	}
	if step.GoBinary != nil {
		return b.BuildModBin(dir, step.GoBinary)
	}
	if step.NodeJS != nil {
		return b.BuildNodeJS(dir, step.NodeJS)
	}
	if step.Dockers != nil {
		return buildDockers(b, dir, step.Dockers)
	}
	if step.DockerPull != nil {
		opt := &elsa.DockerPullOptions{}
		return b.PullDockers(dir, step.DockerPull, opt)
	}
	return nil
}

func buildTarget(
	b *elsa.Builder, targets *targets, name string, opts *buildOptions,
) error {
	step, ok := targets.m[name]
	if ok {
		return buildStep(b, step, opts)
	}
	if name == "base" {
		return b.BuildBase()
	}
	if name == "nodejs" {
		return b.BuildBaseNodeJS()
	}
	return errcode.NotFoundf("not found")
}

func cmdBuild(args []string) error {
	opts := &buildOptions{
		DockerPull: &elsa.DockerPullOptions{},
	}

	flags := cmdFlags.New()
	config := new(elsa.Config)
	declareBuildFlags(flags, config)
	flags.BoolVar(
		&opts.DockerPull.Update,
		"docker_pull_update",
		false,
		"if update when running docker pull",
	)
	args = flags.ParseArgs(args)

	b := elsa.NewBuilder(config)

	build, err := elsa.ReadBuild(buildFile)
	if err != nil {
		return errcode.Annotate(err, "read build")
	}
	ts := newTargets(build)

	if len(args) == 0 {
		for _, step := range ts.steps {
			if err := buildStep(b, step, opts); err != nil {
				return errcode.Annotatef(err, "build %q", step.Name)
			}
		}
	} else {
		for _, target := range args {
			if err := buildTarget(b, ts, target, opts); err != nil {
				return errcode.Annotatef(err, "build %q", target)
			}
		}
	}
	return nil
}
