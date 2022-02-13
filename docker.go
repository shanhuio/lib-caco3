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

package caco3

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/tarutil"
	"shanhu.io/virgo/dock"
)

type dockerImageOptions struct {
	tags     []string
	saveName bool
}

func buildDockerImage(
	env *env, name string, files *tarutil.Stream,
	options *dockerImageOptions,
) error {
	c := env.docker()

	dockName := env.dockerName(name)
	if err := dock.BuildImageStream(c, dockName, files); err != nil {
		return errcode.Annotate(err, "build image")
	}

	log.Printf("saving docker %s", name)
	out, err := env.prepareOut(name + ".tgz")
	if err != nil {
		return errcode.Annotate(err, "prepare output")
	}

	var saveFrom string
	if options != nil && options.saveName {
		saveFrom = dockName
	} else {
		info, err := dock.InspectImage(c, dockName)
		if err != nil {
			return errcode.Annotate(err, "inspect built image")
		}
		saveFrom = info.ID
	}

	if err := dock.SaveImageGz(c, saveFrom, out); err != nil {
		return errcode.Annotate(err, "save output")
	}

	if options != nil {
		for _, t := range options.tags {
			repo, tag := dock.ParseImageTag(t)
			if tag == "" {
				tag = "latest"
			}
			log.Printf("tag as %s:%s", repo, tag)
			if err := dock.TagImage(c, dockName, repo, tag); err != nil {
				return errcode.Annotatef(err, "tag %q", t)
			}
		}
	}

	return nil
}

type docker struct {
	env *env
}

func newDocker(env *env) *docker {
	return &docker{env: env}
}

type dockerFileInput struct {
	tags   []string
	ins    []string
	inZips []string
}

type dockerFileHashTag struct {
	prefix string
	lines  []string
}

func newDockerFileHashTag(tag string) *dockerFileHashTag {
	return &dockerFileHashTag{
		prefix: fmt.Sprintf("#%s ", tag),
	}
}

func (h *dockerFileHashTag) match(line string) bool {
	if !strings.HasPrefix(line, h.prefix) {
		return false
	}
	in := strings.TrimSpace(strings.TrimPrefix(line, h.prefix))
	h.lines = append(h.lines, in)
	return true
}

func (h *dockerFileHashTag) hits() []string {
	return h.lines
}

func parseDockerFileInput(df string) *dockerFileInput {
	tags := newDockerFileHashTag("tag")
	ins := newDockerFileHashTag("in")
	inZips := newDockerFileHashTag("inzip")
	all := []*dockerFileHashTag{tags, ins, inZips}

	s := bufio.NewScanner(strings.NewReader(df))
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}

		matched := false
		for _, t := range all {
			if t.match(line) {
				matched = true
				break
			}
		}
		if !matched { // Hash tag matches must be at file header.
			break
		}
	}

	return &dockerFileInput{
		tags:   tags.hits(),
		ins:    ins.hits(),
		inZips: inZips.hits(),
	}
}

func (d *docker) parseInput(dir, in string) (string, error) {
	var p string
	if strings.HasPrefix(in, "//") {
		p = strings.TrimPrefix(in, "//")
		return d.env.out(p), nil
	}
	if strings.HasPrefix(in, "./") || strings.HasPrefix(in, "../") {
		p = path.Clean(path.Join(dir, in))
		return d.env.src(p), nil
	}

	u, err := url.Parse(in)
	if err != nil {
		return "", err
	}
	if u.Scheme == "src" {
		p = d.env.src(strings.TrimPrefix(u.Path, "/"))
	} else if u.Scheme == "out" {
		p = d.env.out(strings.TrimPrefix(u.Path, "/"))
	}
	return "", errcode.InvalidArgf("unsupported scheme: %q", in)
}

func (d *docker) build(dir string, saveName bool) error {
	srcDir := d.env.src(dir)
	const dockerFileName = "Dockerfile"
	bs, err := ioutil.ReadFile(filepath.Join(srcDir, dockerFileName))
	if err != nil {
		return errcode.Annotate(err, "read Dockerfile")
	}
	df := string(bs)

	in := parseDockerFileInput(df)
	ts := dock.NewTarStream(df)
	for _, in := range in.ins {
		p, err := d.parseInput(dir, in)
		if err != nil {
			return errcode.Annotatef(err, "parse input: %q", in)
		}

		stat, err := os.Stat(p)
		if err != nil {
			return errcode.Annotatef(err, "check file %q", p)
		}
		mode := stat.Mode()
		if !mode.IsRegular() {
			return errcode.InvalidArgf("%q is not a regular file", p)
		}
		ts.AddFile(filepath.Base(p), tarutil.ModeMeta(int64(mode)&0777), p)
	}

	for _, z := range in.inZips {
		p, err := d.parseInput(dir, z)
		if err != nil {
			return errcode.Annotatef(err, "parse zip input: %q", z)
		}
		ts.AddZipFile(p)
	}

	opts := &dockerImageOptions{tags: in.tags, saveName: saveName}
	if err := buildDockerImage(d.env, dir, ts, opts); err != nil {
		return errcode.Annotatef(err, "build docker %q", dir)
	}

	return nil
}
