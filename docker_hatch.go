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
