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

const nodejsDockerfile = `
FROM cr.shanhu.io/base/alpine
MAINTAINER Shanhu Tech Inc.

RUN apk add --update \
	nodejs npm make git openssh zip
RUN npm install -g npm
RUN npm install -g typescript esbuild-linux-64 less

RUN mkdir /usr/local/idle
COPY idle.js /usr/local/idle/idle.js

CMD ["/usr/bin/node", "/usr/local/idle/idle.js"]
`

const idleJS = `
function f() { setTimeout(f, 60*60*1000) }
function exit() {
	console.log('exiting')
	process.exit()
}
process.on('SIGTERM', exit)
process.on('SIGINT', exit)
f()`

var dockerNodejs = &baseDocker{
	name: "base/nodejs",
	build: func(env *env, name string) error {
		ts := dock.NewTarStream(nodejsDockerfile)
		ts.AddString("idle.js", tarutil.ModeMeta(0600), idleJS)
		return buildDockerImage(env, name, nil, ts)
	},
}
