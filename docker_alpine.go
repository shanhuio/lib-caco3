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
	"shanhu.io/virgo/dock"
)

const alpineDockerfile = `
FROM alpine
MAINTAINER Shanhu Tech Inc.

RUN apk update
RUN apk add ca-certificates && update-ca-certificates
`

var dockerAlpine = &baseDocker{
	name: "base/alpine",
	build: func(env *env, name string) error {
		ts := dock.NewTarStream(alpineDockerfile)
		return buildDockerImage(env, name, nil, ts)
	},
}
