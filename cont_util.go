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
	"path"

	"shanhu.io/virgo/dock"
)

func exitError(exit int) error {
	if exit == 0 {
		return nil
	}
	return fmt.Errorf("exit with code: %d", exit)
}

func execError(ret int, err error) error {
	if err != nil {
		return err
	}
	return exitError(ret)
}

func contExec(cont *dock.Cont, args []string) error {
	return execError(cont.ExecArgs(args))
}

func linuxPathJoin(parts ...string) string {
	return path.Join(parts...)
}
