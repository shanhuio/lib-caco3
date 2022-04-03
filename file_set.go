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
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/jsonutil"
	"shanhu.io/misc/strutil"
)

func listAllFiles(dir string) ([]string, error) {
	var files []string
	walk := func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() { // Ignore all directories.
			return nil
		}
		files = append(files, p)
		return nil
	}

	if err := filepath.WalkDir(dir, walk); err != nil {
		return nil, err
	}
	return files, nil
}

type fileSet struct {
	name     string
	files    []string
	includes []string
	out      string // Output file list.
	rule     *FileSet
}

func fileSetOut(name string) string { return name + ".fileset" }

func newFileSet(env *env, p string, r *FileSet) (*fileSet, error) {
	name := makeRelPath(p, r.Name)

	m := make(map[string]bool)
	for _, f := range r.Files {
		m[makePath(p, f)] = true
	}

	var ignores []string
	for _, i := range r.Ignore {
		ignores = append(ignores, makeRelPath(p, i))
	}

	bads := make(map[string]bool)
	ignore := func(name string) bool {
		for _, i := range ignores {
			matched, err := path.Match(i, name)
			if err != nil {
				if !bads[i] {
					log.Printf("bad ignore pattern: %q: %s", i, err)
				}
				bads[i] = true // report each bad ignore pattern once
				continue
			}
			if matched {
				return true
			}
		}
		return false
	}

	for _, sel := range r.Select {
		var matches []string
		// TODO(h8liu): this can be done better.
		if strings.HasSuffix(sel, "/**") {
			dir := env.src(makeRelPath(p, strings.TrimSuffix(sel, "/**")))
			// TODO(h8liu): should also ignore early here when listing.
			files, err := listAllFiles(dir)
			if err != nil {
				return nil, errcode.Annotatef(err, "list all files %q", sel)
			}
			matches = files
		} else {
			glob, err := filepath.Glob(env.src(makeRelPath(p, sel)))
			if err != nil {
				return nil, errcode.Annotatef(err, "glob %q", sel)
			}
			matches = glob
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("%q select no files", sel)
		}

		for _, match := range matches {
			rel, err := filepath.Rel(env.srcDir, match)
			if err != nil {
				return nil, errcode.Annotatef(
					err, "get relative path for %q", match,
				)
			}
			if ignore(rel) {
				continue
			}
			name := filepath.ToSlash(rel)
			m[name] = true
		}
	}

	return &fileSet{
		name:     name,
		files:    strutil.SortedList(m),
		includes: r.Include,
		rule:     r,
		out:      fileSetOut(name),
	}, nil
}

func (fs *fileSet) meta(env *env) (*buildRuleMeta, error) {
	d, err := makeDigest(ruleFileSet, fs.name, fs.rule)
	if err != nil {
		return nil, errcode.Annotate(err, "digest")
	}

	var deps []string
	deps = append(deps, fs.files...)
	deps = append(deps, fs.includes...)

	return &buildRuleMeta{
		name:   fs.name,
		deps:   deps,
		outs:   []string{fs.out},
		digest: d,
	}, nil
}

func referenceFileSet(env *env, name string) (string, error) {
	if t := env.nodeType(name); t == nodeRule {
		if rt := env.ruleType(name); rt != ruleFileSet {
			return "", fmt.Errorf("not a file set, but %q", rt)
		}
		return fileSetOut(name), nil
	} else if t != nodeOut {
		return "", fmt.Errorf("not a file set, but %q", t)
	}
	return name, nil
}

func (fs *fileSet) build(env *env, opts *buildOpts) error {
	m := make(map[string]*fileStat)
	add := func(s *fileStat) {
		// TODO(h8liu): check if files change?
		if _, ok := m[s.Name]; !ok {
			m[s.Name] = s
		}
	}

	for _, f := range fs.files {
		t := env.nodeType(f)
		switch t {
		case "":
			return fmt.Errorf("file %q not found", f)
		case nodeSrc:
			s, err := newSrcFileStat(env, f)
			if err != nil {
				return errcode.Annotatef(err, "file stat %q", f)
			}
			add(s)
		case nodeOut:
			s, err := newOutFileStat(env, f)
			if err != nil {
				return errcode.Annotatef(err, "out file stat %q", f)
			}
			add(s)
		default:
			return fmt.Errorf("unsupported file type %q", t)
		}
	}

	for _, inc := range fs.includes {
		ref, err := referenceFileSet(env, inc)
		if err != nil {
			return errcode.Annotatef(err, "include %q", inc)
		}

		var list []*fileStat
		if err := jsonutil.ReadFile(env.out(ref), &list); err != nil {
			return errcode.Annotatef(err, "read file set %q", inc)
		}
		for _, entry := range list {
			add(entry)
		}
	}

	var names []string
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)

	var list []*fileStat
	for _, name := range names {
		list = append(list, m[name])
	}
	out, err := env.prepareOut(fs.out)
	if err != nil {
		return errcode.Annotate(err, "prepare output")
	}
	return jsonutil.WriteFile(out, list)
}
