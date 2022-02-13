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

package caco3

import (
	"shanhu.io/misc/tarutil"
	"shanhu.io/virgo/dock"
)

const forgeDockerfile = `
FROM cr.shanhu.io/base/golang
MAINTAINER Shanhu Tech Inc.

RUN mkdir /usr/local/idle
COPY idle.go /usr/local/idle/idle.go
RUN cd /usr/local/idle && /usr/local/go/bin/go build idle.go
RUN rm /usr/local/idle/idle.go

CMD ["/usr/local/idle/idle"]
`

var dockerForge = &baseDocker{
	name: "base/forge",
	build: func(env *env, name string) error {
		ts := dock.NewTarStream(forgeDockerfile)
		ts.AddString("idle.go", tarutil.ModeMeta(0600), idleGo)
		return buildDockerImage(env, name, ts, nil)
	},
}
