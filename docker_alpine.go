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
