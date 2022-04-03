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
	"fmt"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/jsonutil"
	"shanhu.io/virgo/dock"
)

type dockerPull struct {
	name    string
	rule    *DockerPull
	repoTag string
	out     string
}

func dockerPullOut(name string) string { return name + ".dockersum" }

func newDockerPull(env *env, p string, r *DockerPull) (*dockerPull, error) {
	name := makeRelPath(p, r.Name)
	repoTag, err := env.nameToRepoTag(name)
	if err != nil {
		return nil, errcode.Annotate(err, "invalid docker pull name")
	}
	return &dockerPull{
		name:    name,
		rule:    r,
		repoTag: repoTag,
		out:     dockerPullOut(name),
	}, nil
}

func (p *dockerPull) pull(env *env) (*dockerSum, error) {
	r := p.rule

	repo, tag := parseRepoTag(p.repoTag)
	srcRepo, srcTag := repo, tag

	if r.Pull != "" {
		srcRepo, srcTag = parseRepoTag(r.Pull)
	}

	digest := r.Digest

	from := repoTag(srcRepo, srcTag)
	pullTag := srcTag

	if digest != "" {
		from = fmt.Sprintf("%s@%s", srcRepo, digest)
		pullTag = digest
	}

	if err := dock.PullImage(env.dock, srcRepo, pullTag); err != nil {
		return nil, errcode.Annotate(err, "pull image")
	}
	if err := dock.TagImage(env.dock, from, srcRepo, srcTag); err != nil {
		return nil, errcode.Annotate(err, "tag image as source")
	}
	if !(repo == srcRepo && tag == srcTag) {
		if err := dock.TagImage(env.dock, from, repo, tag); err != nil {
			return nil, errcode.Annotate(err, "re-tag output image")
		}
	}
	out := repoTag(repo, tag)
	info, err := dock.InspectImage(env.dock, out)
	if err != nil {
		return nil, errcode.Annotate(err, "inspect image")
	}

	sum := newDockerSum(info, srcRepo, digest)
	if sum.Digest == "" {
		return nil, fmt.Errorf("no digest found for %q", out)
	}
	if digest != "" && sum.Digest != digest {
		return nil, fmt.Errorf(
			"digest mismatch, got %q, want %q", sum.Digest, digest,
		)
	}
	return sum, nil
}

func (p *dockerPull) build(env *env, opts *buildOpts) error {
	sums, err := p.pull(env)
	if err != nil {
		return err
	}
	out, err := env.prepareOut(p.out)
	if err != nil {
		return errcode.Annotate(err, "prepare sums output")
	}
	if err := jsonutil.WriteFile(out, sums); err != nil {
		return errcode.Annotate(err, "write sums")
	}
	return nil
}

func (p *dockerPull) meta(env *env) (*buildRuleMeta, error) {
	digest := ""
	if p.rule.Digest != "" {
		d, err := makeDigest(ruleDockerPull, p.name, p.rule)
		if err != nil {
			return nil, errcode.Annotate(err, "digest")
		}
		digest = d
	}

	return &buildRuleMeta{
		name:   p.name,
		outs:   []string{p.out},
		digest: digest,
	}, nil
}
