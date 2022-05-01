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

package caco3bin

import (
	"log"
	"os"
	"path"

	"shanhu.io/caco3"
	"shanhu.io/misc/errcode"
	"shanhu.io/text/lexing"
)

type targets struct {
	m     map[string]*caco3.BuildStep
	steps []*caco3.BuildStep
}

func newTargets(ws *caco3.Workspace) *targets {
	m := make(map[string]*caco3.BuildStep)
	for _, step := range ws.Steps {
		m[step.Name] = step
	}
	return &targets{
		m:     m,
		steps: ws.Steps,
	}
}

func buildDockers(
	b *caco3.Builder, dir string, docks []string, saveName bool,
) error {
	for _, d := range docks {
		if err := b.BuildDocker(path.Join(dir, d), saveName); err != nil {
			return errcode.Annotatef(err, "build docker %q", d)
		}
	}
	return nil
}

type buildOptions struct {
	DockerPull     *caco3.DockerPullOptions
	DockerSaveName bool
}

func buildStep(
	b *caco3.Builder, step *caco3.BuildStep, opts *buildOptions,
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
		return buildDockers(b, dir, step.Dockers, opts.DockerSaveName)
	}
	if step.DockerPull != nil {
		return b.PullDockers(dir, step.DockerPull, opts.DockerPull)
	}
	return nil
}

func buildTarget(
	b *caco3.Builder, targets *targets, name string, opts *buildOptions,
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
		DockerPull: &caco3.DockerPullOptions{},
	}

	flags := cmdFlags.New()
	config := new(caco3.Config)
	declareBuildFlags(flags, config)
	flags.BoolVar(
		&opts.DockerPull.Update,
		"docker_pull_update",
		false,
		"if update when running docker pull",
	)
	args = flags.ParseArgs(args)

	wd, err := os.Getwd()
	if err != nil {
		return errcode.Annotate(err, "get work dir")
	}

	b := caco3.NewBuilder(wd, config)
	ws, errs := b.ReadWorkspace()
	if errs != nil {
		lexing.FprintErrs(os.Stderr, errs, wd)
		return errcode.InvalidArgf("read build got %d errors", len(errs))
	}

	opts.DockerSaveName = ws.DockerSaveName
	ts := newTargets(ws)
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

func cmdBuild2(args []string) error {
	flags := cmdFlags.New()
	config := new(caco3.Config)
	declareBuildFlags(flags, config)
	args = flags.ParseArgs(args)

	wd, err := os.Getwd()
	if err != nil {
		return errcode.Annotate(err, "get work dir")
	}

	b := caco3.NewBuilder(wd, config)
	if errs := b.Build(args); errs != nil {
		lexing.FprintErrs(os.Stderr, errs, wd)
		return errcode.InvalidArgf("build got %d errors", len(errs))
	}

	return nil
}
