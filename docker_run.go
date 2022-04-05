package caco3

import (
	"log"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/strutil"
	"shanhu.io/misc/tarutil"
	"shanhu.io/virgo/dock"
)

type dockerRun struct {
	name   string
	rule   *DockerRun
	image  string
	ins    []string
	inMap  map[string]string
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

	depsMap := make(map[string]bool)
	for _, d := range r.Deps {
		depsMap[makePath(p, d)] = true
	}

	var ins []string
	inMap := make(map[string]string)
	for f, v := range r.Input {
		inPath := makePath(p, f)
		ins = append(ins, inPath)
		inMap[inPath] = v
		depsMap[inPath] = true
	}
	deps = append(deps, strutil.SortedList(depsMap)...)

	var outs []string
	outMap := make(map[string]string)
	for f, v := range r.Output {
		outPath := makeRelPath(p, f)
		outs = append(outs, outPath)
		outMap[outPath] = v
	}
	outs = strutil.SortedList(strutil.MakeSet(outs))

	return &dockerRun{
		name:   name,
		rule:   r,
		image:  image,
		ins:    ins,
		inMap:  inMap,
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

	if len(r.ins) > 0 {
		ts := tarutil.NewStream()
		for _, in := range r.ins {
			var f string
			switch typ := env.nodeType(in); typ {
			case "":
				return errcode.Internalf("input %q not found", in)
			case nodeSrc:
				f = env.src(in)
			case nodeOut:
				f = env.out(in)
			default:
				return errcode.Internalf("unknown type %q", typ)
			}

			dest := r.inMap[in]
			ts.AddFile(dest, new(tarutil.Meta), f)
		}

		if err := dock.CopyInTarStream(cont, ts, "/"); err != nil {
			return errcode.Annotate(err, "copy input")
		}
	}

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
