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

package elsa

import (
	"fmt"
	"runtime"

	"shanhu.io/virgo/dock"
)

type golang struct {
	env        *env
	hatchImage string
}

func newGolang(env *env, hatchImage string) *golang {
	return &golang{
		env:        env,
		hatchImage: hatchImage,
	}
}

func (g *golang) build(d *dock.Client, source *golangSource) error {
	srcTgz := g.env.out("go-src.tar.gz")
	if err := source.downloadTo(srcTgz); err != nil {
		return err
	}

	c, err := dock.CreateCont(d, g.hatchImage, nil)
	if err != nil {
		return err
	}
	defer c.ForceRemove()

	if err := c.Start(); err != nil {
		return err
	}

	// install packages
	if err := dock.RunTasks(c, []string{
		"apk update",
		"apk add --no-cache ca-certificates",
		"apk add --no-cache bash gcc musl-dev openssl go",
	}); err != nil {
		return err
	}

	if err := dock.CopyInTarGz(c, srcTgz, "/usr/local"); err != nil {
		return err
	}

	// compile the golang source
	exit, err := c.ExecWithSetup(&dock.ExecSetup{
		Cmd: []string{"/bin/bash", "make.bash"},
		Env: []string{
			"GOOS=linux",
			"GOARCH=" + runtime.GOARCH,
			"GOHOSTOS=linux",
			"GOHOSTARCH=" + runtime.GOARCH,
			"GOROOT_BOOTSTRAP=/usr/lib/go",
		},
		WorkingDir: "/usr/local/go/src",
	})
	if err != nil {
		return err
	}
	if exit != 0 {
		return fmt.Errorf("exit value: %d", exit)
	}

	// clean up useless stuff
	if err := dock.RunTask(
		c, "rm -rf /usr/local/go/pkg/bootstrap /usr/local/go/pkg/obj",
	); err != nil {
		return err
	}

	output := g.env.out("go.tar.gz")
	if err := dock.CopyOutTarGz(c, "/usr/local/go", output); err != nil {
		return err
	}

	return c.Drop()
}
