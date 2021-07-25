package elsa

import (
	"shanhu.io/misc/errcode"
	"shanhu.io/misc/jsonx"
	"shanhu.io/misc/tarutil"
	"shanhu.io/virgo/dock"
)

const golangDockerfile = `
FROM cr.shanhu.io/base/alpine
MAINTAINER Shanhu Tech Inc.

RUN apk add --update \
	git subversion mercurial \
	gcc g++ musl-dev make openssh
ADD go.tar.gz /usr/local
ENV PATH /usr/local/go/bin:/usr/sbin:/usr/bin:/sbin:/bin
RUN mkdir /go
ENV GOPATH /go
WORKDIR /go

CMD ["/usr/local/go/bin/go", "version"]
`

var dockerGolang = &baseDocker{
	name: "base/golang",
	build: func(env *env, name string) error {
		src := new(golangSource)
		srcFile := env.src(env.goVersion)
		if err := jsonx.ReadFile(srcFile, src); err != nil {
			return errcode.Annotate(err, "read golang source config")
		}
		golang := newGolang(env, env.dockerName(dockerHatch.name))
		if err := golang.build(env.docker(), src); err != nil {
			return err
		}

		ts := dock.NewTarStream(golangDockerfile)
		ts.AddFile("go.tar.gz", tarutil.ModeMeta(0600), env.out("go.tar.gz"))
		return buildDockerImage(env, name, nil, ts)
	},
}
