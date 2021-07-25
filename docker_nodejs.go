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
