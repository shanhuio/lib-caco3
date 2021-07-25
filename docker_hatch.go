// Copyright (C) 2021  Shanhu Tech Inc.
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
	"shanhu.io/misc/tarutil"
	"shanhu.io/virgo/dock"
)

const hatchDockerfile = `
FROM cr.shanhu.io/base/alpine
MAINTAINER Shanhu Tech Inc.

RUN apk update && apk add musl-dev go
RUN mkdir /hatch
COPY idle.go /hatch/idle.go
RUN cd /hatch && go build idle.go
WORKDIR /hatch
CMD ["/hatch/idle"]
`

// This idle program is required to perform as an "sleep infinity", but
// unlike sleep, it can receive signals and be stopped at any time.
const idleGo = `package main
import "time"
func main() { for { time.Sleep(time.Hour) } }
`

var dockerHatch = &baseDocker{
	name: "base/hatch",
	build: func(env *env, name string) error {
		ts := dock.NewTarStream(hatchDockerfile)
		ts.AddString("idle.go", tarutil.ModeMeta(0600), idleGo)
		return buildDockerImage(env, name, nil, ts)
	},
}
