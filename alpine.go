package elsa

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"shanhu.io/misc/errcode"
	"shanhu.io/virgo/dock"
)

type alpine struct {
	env   *env
	image string // Docker image to the forge, normally "shanhu/forge".
}

func newAlpine(env *env, image string) *alpine {
	return &alpine{
		env:   env,
		image: image,
	}
}

func (a *alpine) binPath(pkg string) string {
	return a.env.out("alpine-bin", pkg)
}

func (a *alpine) buildMod(dir string, pkgs []string) error {
	absSrc, err := filepath.Abs(a.env.src())
	if err != nil {
		return errcode.Annotate(err, "get absolute src dir")
	}

	const srcRoot = "/go/elsa-src"

	contConfig := &dock.ContConfig{
		Mounts: []*dock.ContMount{{
			Host:     absSrc,
			Cont:     srcRoot,
			ReadOnly: true,
		}},
	}
	client := a.env.docker()
	cont, err := dock.CreateCont(client, a.image, contConfig)
	if err != nil {
		return err
	}
	defer cont.Drop()

	if err := cont.Start(); err != nil {
		return err
	}

	const outRoot = "/go/elsa-out"

	if err := contExec(cont, []string{"mkdir", "-p", outRoot}); err != nil {
		return errcode.Annotate(err, "make output root")
	}

	envVars := []string{
		"GOPATH=/go",
		"GO111MODULE=on",
		"GOPRIVATE=shanhu.io",
	}

	if a.env.sshKnownHosts != "" {
		f := linuxPathJoin(srcRoot, a.env.sshKnownHosts)
		const hostsFile = "/go/ssh_known_hosts"
		if err := contExec(cont, []string{"cp", f, hostsFile}); err != nil {
			return errcode.Annotate(err, "setup known hosts")
		}
		envVars = append(envVars, fmt.Sprintf(
			"GIT_SSH_COMMAND=ssh -o UserKnownHostsFile='%s'", hostsFile,
		))
	}
	workDir := linuxPathJoin(srcRoot, dir)

	for _, pkg := range pkgs {
		log.Printf("bin %s", pkg)

		base := path.Base(pkg)
		pkgDir := path.Clean(path.Dir(pkg))

		binDir := linuxPathJoin(outRoot, dir, pkgDir)

		for _, args := range [][]string{
			{"mkdir", "-p", binDir},
			{
				"go", "build",
				"-ldflags=-s -w", "-trimpath",
				"-o", linuxPathJoin(binDir, base),
				pkg,
			},
		} {
			exit, err := cont.ExecWithSetup(&dock.ExecSetup{
				Cmd:        args,
				Env:        envVars,
				WorkingDir: workDir,
			})
			if err != nil {
				return err
			}
			if exit != 0 {
				return exitError(exit)
			}
		}
	}

	outDest := filepath.Dir(a.env.out(dir))
	if err := os.MkdirAll(outDest, 0700); err != nil {
		return err
	}
	return cont.CopyOut(linuxPathJoin(outRoot, dir), outDest)
}

func (a *alpine) build(pkgs []string) error {
	for _, pkg := range pkgs {
		if strings.HasSuffix(pkg, "/") {
			return errcode.InvalidArgf("%q is not a valid package", pkg)
		}
	}

	log.Println("setup alpine builder docker instance")

	absSrc, err := filepath.Abs(a.env.goSrc())
	if err != nil {
		return errcode.Annotate(err, "get absolute gopath/src")
	}

	contConfig := &dock.ContConfig{
		Mounts: []*dock.ContMount{{
			Host:     absSrc,
			Cont:     "/go/src",
			ReadOnly: true,
		}},
	}
	client := a.env.docker()
	cont, err := dock.CreateCont(client, a.image, contConfig)
	if err != nil {
		return err
	}
	defer cont.Drop()

	if err := cont.Start(); err != nil {
		return err
	}

	const alpineBin = "/go/alpine-bin"

	if err := contExec(cont, []string{"mkdir", "-p", alpineBin}); err != nil {
		return err
	}

	envVars := []string{
		"GOPATH=/go",
		"GO111MODULE=off",
	}

	for _, pkg := range pkgs {
		log.Printf("bin %s", pkg)

		base := path.Base(pkg)
		dir := path.Dir(pkg)

		for _, args := range [][]string{
			{"rm", "-rf", "/go/bin"},
			{"go", "install", "-ldflags=-s -w", "-trimpath", pkg},
			{"mkdir", "-p", linuxPathJoin(alpineBin, dir)},
			{
				"mv",
				linuxPathJoin("/go/bin/", base),
				linuxPathJoin(alpineBin, pkg),
			},
		} {
			if err := execError(cont.ExecWithSetup(&dock.ExecSetup{
				Cmd:        args,
				Env:        envVars,
				WorkingDir: "/go",
			})); err != nil {
				return errcode.Annotatef(err, "%q", args)
			}
		}
	}

	output := a.env.out()
	if err := os.MkdirAll(output, 0700); err != nil {
		return err
	}
	return cont.CopyOut(alpineBin, output)
}
