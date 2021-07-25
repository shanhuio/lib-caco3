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
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/idutil"
	"shanhu.io/misc/osutil"
)

func currentGitCommit(dir string) (string, error) {
	branches, err := runCmdOutput(dir, "git", "branch")
	if err != nil {
		return "", errcode.Annotate(err, "list branches")
	}
	if len(bytes.TrimSpace(branches)) == 0 {
		return "", nil
	}

	ret, err := runCmdOutput(
		dir, "git", "show", "HEAD", "-s", "--format=%H",
	)
	if err != nil {
		return "", errcode.Annotate(err, "get HEAD commit")
	}
	return strings.TrimSpace(string(ret)), nil
}

func gitSync(name, dir, remote, commit string) (bool, error) {
	if commit == "" {
		latest, err := runCmdOutput(dir, "git", "ls-remote", remote, "HEAD")
		if err != nil {
			return false, errcode.Annotate(err, "git ls-remote")
		}
		line := strings.TrimSpace(string(latest))
		fields := strings.Fields(line)
		if len(fields) == 0 {
			return false, errcode.Internalf("bad remote commit: %q", line)
		}
		commit = fields[0]
	}

	gitDir := filepath.Join(dir, ".git")
	exist, err := osutil.IsDir(gitDir)
	if err != nil {
		return false, errcode.Annotate(err, "check git dir")
	}

	const stashBranch = "elsa"

	if !exist {
		if err := runCmd(dir, "git", "init", "-q"); err != nil {
			return false, errcode.Annotate(err, "git init")
		}
		if err := runCmd(
			dir, "git", "remote", "add", "origin", remote,
		); err != nil {
			return false, errcode.Annotate(err, "git add remote")
		}

		log.Printf(
			"[new %s] %s\n", idutil.Short(commit), name,
		)
	} else {
		cur, err := currentGitCommit(dir)
		if err != nil {
			return false, errcode.Annotate(err, "get current comment")
		}
		if cur == commit {
			return false, nil
		}

		if cur != "" {
			hasCommit, err := callCmd(
				dir, "git", "cat-file", "-e", commit,
			)
			if err != nil {
				return false, errcode.Annotate(err, "git check commit")
			}
			if hasCommit {
				isAncestor, err := callCmd(
					dir, "git", "merge-base", "--is-ancestor", commit, cur,
				)
				if err != nil {
					return false, errcode.Annotate(err, "git merge check")
				}
				if isAncestor {
					// merge will be a noop, just update stash branch.
					return false, runCmd(
						dir, "git", "branch", "-q", "-f", stashBranch, commit,
					)
				}
			}

			log.Printf(
				"[%s..%s] %s\n",
				idutil.Short(cur), idutil.Short(commit), name,
			)
		} else {
			log.Printf(
				"[new %s] %s\n", idutil.Short(commit), name,
			)
		}
	}

	// fetch to the stash branch and then merge.
	if err := runCmd(
		dir, "git", "fetch", "-q", remote, "HEAD",
	); err != nil {
		return false, errcode.Annotate(err, "git fetch")
	}
	if err := runCmd(
		dir, "git", "branch", "-q", "-f", stashBranch, commit,
	); err != nil {
		return false, errcode.Annotate(err, "git branch stash")
	}
	if err := runCmd(
		dir, "git", "merge", "-q", stashBranch,
	); err != nil {
		return false, errcode.Annotate(err, "git merge stash")
	}

	return true, nil
}

func syncRepos(env *env, b *Build, sums *BuildSums) error {
	var dirs []string
	for dir := range b.Repos {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)

	for _, dir := range dirs {
		git := b.Repos[dir]
		srcDir := env.src(dir)
		if err := os.MkdirAll(srcDir, 0755); err != nil {
			return errcode.Annotatef(err, "make dir for %q", dir)
		}

		commit := ""
		if sums != nil {
			c, ok := sums.RepoCommits[dir]
			if !ok {
				return errcode.InvalidArgf("commit missing for %q", dir)
			}
			commit = c
		}
		if _, err := gitSync(dir, srcDir, git, commit); err != nil {
			return errcode.Annotatef(err, "git sync %q", dir)
		}

	}
	return nil
}
