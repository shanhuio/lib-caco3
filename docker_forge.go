package elsa

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
		return buildDockerImage(env, name, nil, ts)
	},
}
