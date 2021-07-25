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

func buildStep(b *elsa.Builder, step *elsa.BuildStep) error {
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

func buildTarget(b *elsa.Builder, targets *targets, name string) error {
	step, ok := targets.m[name]
	if ok {
		return buildStep(b, step)
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
	flags := cmdFlags.New()
	config := new(elsa.Config)
	declareBuildFlags(flags, config)
	flags.ParseArgs(args)

	b := elsa.NewBuilder(config)

	build, err := elsa.ReadBuild(buildFile)
	if err != nil {
		return errcode.Annotate(err, "read build")
	}
	ts := newTargets(build)

	if len(args) == 0 {
		for _, step := range ts.steps {
			if err := buildStep(b, step); err != nil {
				return errcode.Annotatef(err, "build %q", step.Name)
			}
		}
	} else {
		for _, target := range args {
			if err := buildTarget(b, ts, target); err != nil {
				return errcode.Annotatef(err, "build %q", target)
			}
		}
	}
	return nil
}
