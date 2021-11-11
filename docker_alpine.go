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
	"text/template"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/jsonx"
	"shanhu.io/virgo/dock"
)

const alpineDockerfileContent = `
FROM alpine:{{.Version}}
MAINTAINER Shanhu Tech Inc.

RUN apk update
RUN apk add ca-certificates && update-ca-certificates
`

var alpineDockerfileTmpl = template.Must(
	template.New("_").Parse(alpineDockerfileContent),
)

func makeAlpineDockerFile(ver string) (string, error) {
	buf := new(bytes.Buffer)
	dat := struct {
		Version string
	}{
		Version: ver,
	}
	if err := alpineDockerfileTmpl.Execute(buf, &dat); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var dockerAlpine = &baseDocker{
	name: "base/alpine",
	build: func(env *env, name string) error {
		var version struct {
			Version string
		}
		versionFile := env.src("base/alpine.jsonx")
		if err := jsonx.ReadFile(versionFile, &version); err != nil {
			return errcode.Annotate(err, "read alpine version")
		}

		dockerFile, err := makeAlpineDockerFile(version.Version)
		if err != nil {
			return errcode.Annotate(err, "make Dockerfile")
		}

		ts := dock.NewTarStream(dockerFile)
		return buildDockerImage(env, name, ts, nil)
	},
}
