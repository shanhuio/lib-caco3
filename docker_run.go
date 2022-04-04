package caco3

import (
	"log"
	"sort"

	"shanhu.io/misc/errcode"
	"shanhu.io/virgo/dock"
)

type dockerRun struct {
	name   string
	rule   *DockerRun
	image  string
	deps   []string
	outs   []string
	outMap map[string]string
	envs   map[string]string
}

func newDockerRun(env *env, p string, r *DockerRun) *dockerRun {
	name := makeRelPath(p, r.Name)

	image := makePath(p, r.Image)
	var deps []string
	deps = append(deps, image)
	for _, d := range r.Deps {
		deps = append(deps, makePath(p, d))
	}
	sort.Strings(deps)

	var outs []string
	outMap := make(map[string]string)
	for f, v := range r.Output {
		outPath := makeRelPath(p, f)
		outs = append(outs, outPath)
		outMap[outPath] = v
	}
	sort.Strings(outs)

	return &dockerRun{
		name:   name,
		rule:   r,
		image:  image,
		deps:   deps,
		outs:   outs,
		outMap: outMap,
		envs:   makeDockerVars(r.Envs),
	}
}

func (r *dockerRun) meta(env *env) (*buildRuleMeta, error) {
	dat := struct {
		Rule *DockerRun
		Envs map[string]string `json:",omitempty"`
	}{
		Rule: r.rule,
		Envs: r.envs,
	}
	digest, err := makeDigest(ruleDockerRun, r.name, &dat)
	if err != nil {
		return nil, errcode.Annotate(err, "digest")
	}

	return &buildRuleMeta{
		name:   r.name,
		outs:   r.outs,
		deps:   r.deps,
		digest: digest,
	}, nil
}

func (r *dockerRun) build(env *env, opts *buildOpts) error {
	contConfig := &dock.ContConfig{
		Cmd:     r.rule.Command,
		WorkDir: r.rule.WorkDir,
		Env:     r.envs,
	}

	img, err := env.nameToRepoTag(r.image)
	if err != nil {
		return errcode.Annotate(err, "map image name")
	}

	c := env.dock

	cont, err := dock.CreateCont(c, img, contConfig)
	if err != nil {
		return errcode.Annotate(err, "create container")
	}
	defer cont.Drop()

	if err := cont.Start(); err != nil {
		return errcode.Annotate(err, "start container")
	}
	if err := cont.FollowLogs(opts.log); err != nil {
		return errcode.Annotate(err, "stream logs")
	}

	status, err := cont.Wait(dock.NotRunning)
	if err != nil {
		return errcode.Annotate(err, "wait container")
	}
	for _, out := range r.outs {
		from := r.outMap[out]
		to := out

		f, err := env.prepareOut(to)
		if err != nil {
			return errcode.Annotatef(err, "prepare output: %s", to)
		}

		if err := cont.CopyOutFile(from, f); err != nil {
			if status == 0 {
				return errcode.Annotatef(err, "copy %s", to)
			}
			log.Printf("copy %s: %s", to, err)
		}
	}

	if status != 0 {
		return errcode.Internalf("exit with %d", status)
	}

	return nil
}
