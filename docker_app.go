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
		return buildDockerImage(env, name, nil, ts)
	},
}
