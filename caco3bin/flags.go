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
	"shanhu.io/caco3"
	"shanhu.io/misc/flagutil"
)

var cmdFlags = flagutil.NewFactory("caco3")

func declareBuildFlags(flags *flagutil.FlagSet, c *caco3.Config) {
	flags.StringVar(&c.Root, "root", ".", "root directory")
	flags.StringVar(&c.Src, "src", "src", "source directory")
	flags.StringVar(&c.Out, "out", "out", "output directory")
	flags.StringVar(&c.GoSrc, "gosrc", "", "go language source directory")
	flags.StringVar(&c.DockerRegistry, "cr", "cr.shanhu.io", "docker registry")
	flags.StringVar(
		&c.GoVersion, "goversion", "base/go-src.jsonx",
		"go language version spec file",
	)
	flags.StringVar(
		&c.SSHKnownHosts, "ssh_known_hosts", "base/ssh_known_hosts",
		"ssh known hosts file",
	)
}