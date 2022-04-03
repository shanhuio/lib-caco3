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
	"shanhu.io/misc/subcmd"
)

func cmd() *subcmd.List {
	c := subcmd.New()
	c.Add("build", "builds a target", cmdBuild)
	c.Add("build2", "build rules", cmdBuild2)
	c.Add("sync", "sync source repos", cmdSync)
	return c
}

// Main is the entrance for the elsa binary.
func Main() { cmd().Main() }
