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
	"fmt"
	"log"
	"path"
	"runtime"
	"sort"
	"strings"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/jsonx"
	"shanhu.io/virgo/dock"
)

// DockerSum captures the docker's ID and digest.
type DockerSum struct {
	ID     string
	Digest string
}

func pullDockerImage(
	c *dock.Client, outName, repo, tag, digest string,
) (*DockerSum, error) {
	name := fmt.Sprintf("%s:%s", repo, tag)

	if tag == "" {
		tag = "latest"
	}
	srcTag := tag
	if digest != "" {
		srcTag = digest
	}
	log.Printf("pull docker %s", name)
	if err := dock.PullImage(c, repo, srcTag); err != nil {
		return nil, errcode.Annotate(err, "pull image")
	}
	if srcTag != tag {
		from := fmt.Sprintf("%s@%s", repo, srcTag)
		if err := dock.TagImage(c, from, repo, tag); err != nil {
			return nil, errcode.Annotate(err, "retag image")
		}
	}

	info, err := dock.InspectImage(c, name)
	if err != nil {
		return nil, errcode.Annotate(err, "inspect image")
	}

	sum := &DockerSum{ID: info.ID}
	digestPrefix := repo + "@"
	for _, d := range info.RepoDigests {
		if strings.HasPrefix(d, digestPrefix) {
			sum.Digest = strings.TrimPrefix(d, digestPrefix)
			break
		}
	}
	if sum.Digest == "" {
		return nil, errcode.Internalf("no digest found")
	}
	return sum, nil
}

// DockerPullOptions specifies how the dockers are being pulled.
type DockerPullOptions struct {
	Update     bool
	IgnoreSums bool
}

func pullDockers(
	env *env, dir string, p *DockerPull, opt *DockerPullOptions,
) error {
	imagesFile := env.src(dir, p.Images)
	src := make(map[string]string)
	if err := jsonx.ReadFile(imagesFile, &src); err != nil {
		return errcode.Annotate(err, "read images file")
	}
	if len(src) == 0 {
		return nil
	}
	var names []string
	for name := range src {
		names = append(names, name)
	}
	sort.Strings(names)

	arch := runtime.GOARCH
	sumsFile, ok := p.Sums[arch]
	if !ok {
		return errcode.InvalidArgf("sums file not found for arch %q", arch)
	}
	sumsFile = env.src(dir, sumsFile)

	var sums map[string]*DockerSum
	ignoreSums := opt.IgnoreSums || opt.Update
	if ignoreSums {
		sums := make(map[string]*DockerSum)
		if err := jsonx.ReadFile(sumsFile, &sums); err != nil {
			return errcode.Annotate(err, "read sums file")
		}
	}

	c := dock.NewUnixClient("")

	newSums := make(map[string]*DockerSum)
	for _, name := range names {
		img := src[name]
		repo, tag := dock.ParseImageTag(img)
		if tag == "" {
			tag = "latest"
		}

		var digest string
		if !ignoreSums {
			if sum := sums[img]; sum != nil {
				digest = sum.Digest
			}
		}

		outName := path.Join(dir, name)
		gotSum, err := pullDockerImage(c, outName, repo, tag, digest)
		if err != nil {
			return errcode.Annotatef(err, "pull %s / %s", name, img)
		}

		log.Printf("saving docker %s", outName)
		out, err := env.prepareOut(outName + ".tgz")
		if err != nil {
			return errcode.Annotate(err, "prepare output")
		}
		if err := dock.SaveImageGz(c, gotSum.ID, out); err != nil {
			return errcode.Annotate(err, "save output")
		}

		if opt.Update {
			newSums[img] = gotSum
		}
	}

	if opt.Update {
		if err := jsonx.WriteFile(sumsFile, newSums); err != nil {
			return errcode.Annotate(err, "update new sums")
		}
	}

	return nil
}
