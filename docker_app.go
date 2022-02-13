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
	"shanhu.io/virgo/dock"
)

const appDockerFile = `
FROM cr.shanhu.io/base/alpine
MAINTAINER Shanhu Tech Inc.

RUN adduser -D -u 3000 app
RUN mkdir -p /opt/app
RUN cd /opt/app && mkdir bin var etc lib tmp
RUN chown -R app /opt/app
WORKDIR /opt/app
`

var dockerApp = &baseDocker{
	name: "base/app",
	build: func(env *env, name string) error {
		ts := dock.NewTarStream(appDockerFile)
		return buildDockerImage(env, name, ts, nil)
	},
}
