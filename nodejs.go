package elsa

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"shanhu.io/misc/errcode"
	"shanhu.io/virgo/dock"
)

type nodeJS struct {
	env   *env
	image string
}

func newNodeJS(env *env, image string) *nodeJS {
	return &nodeJS{
		env:   env,
		image: image,
	}
}

func (n *nodeJS) build(dir string, spec *NodeJS) error {
	absSrc, err := filepath.Abs(n.env.src())
	if err != nil {
		return errcode.Annotate(err, "get absolute src dir")
	}

	absOut, err := filepath.Abs(n.env.out())
	if err != nil {
		return errcode.Annotate(err, "get absolute out dir")
	}

	const srcRoot = "/node/elsa-src"
	const depRoot = "/node/elsa-dep"

	contConfg := &dock.ContConfig{
		Mounts: []*dock.ContMount{{
			Host:     absSrc,
			Cont:     srcRoot,
			ReadOnly: true,
		}, {
			Host:     absOut,
			Cont:     depRoot,
			ReadOnly: true,
		}},
	}
	client := n.env.docker()
	cont, err := dock.CreateCont(client, n.image, contConfg)
	if err != nil {
		return errcode.Annotate(err, "create container")
	}
	defer cont.Drop()

	if err := cont.Start(); err != nil {
		return errcode.Annotate(err, "start container")
	}

	const (
		outRoot  = "/node/elsa-out"
		workRoot = "/node/work"
		tmpRoot  = "/node/tmp"
	)

	for _, dir := range []string{
		outRoot, workRoot, tmpRoot,
	} {
		if err := contExec(cont, []string{"mkdir", "-p", dir}); err != nil {
			return errcode.Annotatef(err, "make %q", dir)
		}
	}

	var cmds [][]string

	workDirDir := path.Clean(linuxPathJoin(workRoot, path.Dir(dir)))
	if workDirDir != workRoot {
		cmds = append(cmds, []string{"mkdir", "-p", workDirDir})
	}
	cmds = append(cmds, []string{
		"cp", "-R", linuxPathJoin(srcRoot, dir), workDirDir,
	})

	for _, args := range cmds {
		if err := contExec(cont, args); err != nil {
			return errcode.Annotatef(err, "%q", args)
		}
	}

	var envVars []string
	if n.env.sshKnownHosts != "" {
		f := linuxPathJoin(srcRoot, n.env.sshKnownHosts)
		const hostsFile = "/node/ssh_known_hosts"
		if err := contExec(cont, []string{"cp", f, hostsFile}); err != nil {
			return errcode.Annotate(err, "setup known hosts")
		}
		envVars = append(envVars, fmt.Sprintf(
			"GIT_SSH_COMMAND=ssh -o UserKnownHostsFile='%s'", hostsFile,
		))
	}

	workDir := linuxPathJoin(workRoot, dir)
	for _, args := range [][]string{
		{"npm", "ci"},    // install other deps
		{"make", "dist"}, // make the thing
	} {
		if err := execError(cont.ExecWithSetup(&dock.ExecSetup{
			Cmd:        args,
			Env:        envVars,
			WorkingDir: workDir,
		})); err != nil {
			return errcode.Annotatef(err, "%q", args)
		}
	}

	var outFiles []string
	if spec.Output == nil {
		outFiles = []string{path.Base(dir) + ".tgz"}
	} else {
		outFiles = spec.Output
	}

	if err := contExec(cont, []string{
		"mkdir", "-p", linuxPathJoin(outRoot, dir),
	}); err != nil {
		return errcode.Annotatef(err, "mkdir for output")
	}
	for _, f := range outFiles {
		if err := contExec(cont, []string{
			"mv", linuxPathJoin(workDir, f),
			linuxPathJoin(outRoot, dir, f),
		}); err != nil {
			return errcode.Annotatef(err, "move output %q", f)
		}
	}

	outDest := filepath.Dir(n.env.out(dir))
	if err := os.MkdirAll(outDest, 0700); err != nil {
		return err
	}
	return cont.CopyOut(linuxPathJoin(outRoot, dir), outDest)
}
