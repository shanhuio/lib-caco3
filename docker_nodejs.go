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
	"bytes"
	"runtime"
	"text/template"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/tarutil"
	"shanhu.io/virgo/dock"
)

const nodejsDockerfileContent = `
FROM cr.shanhu.io/base/alpine
MAINTAINER Shanhu Tech Inc.

RUN apk add --update \
	nodejs npm make git openssh zip
RUN npm install -g npm
RUN npm install -g typescript less {{.Esbuild}}

RUN mkdir /usr/local/idle
COPY idle.js /usr/local/idle/idle.js

CMD ["/usr/bin/node", "/usr/local/idle/idle.js"]
`

var nodejsDockerfileTmpl = template.Must(
	template.New("_").Parse(nodejsDockerfileContent),
)

func makeNodeJSDockerFile() (string, error) {
	buf := new(bytes.Buffer)

	var dat struct {
		Esbuild string
	}

	switch runtime.GOARCH {
	case "amd64":
		dat.Esbuild = "esbuild-linux-64"
	case "arm64":
		dat.Esbuild = "esbuild-linux-arm64"
	default:
		dat.Esbuild = "esbuild" // might not work..
	}

	if err := nodejsDockerfileTmpl.Execute(buf, &dat); err != nil {
		return "", err
	}
	return buf.String(), nil
}

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
		dockerFile, err := makeNodeJSDockerFile()
		if err != nil {
			return errcode.Annotate(err, "make DockerFile")
		}
		ts := dock.NewTarStream(dockerFile)
		ts.AddString("idle.js", tarutil.ModeMeta(0600), idleJS)
		return buildDockerImage(env, name, nil, ts)
	},
}
