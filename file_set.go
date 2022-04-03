package caco3

import (
	"fmt"
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"strings"

	"shanhu.io/misc/errcode"
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
	d, err := makeRuleDigest("fileSet", fs.name, fs.rule)
	if err != nil {
		return nil, errcode.Annotate(err, "digest")
	}

	return &buildRuleMeta{
		name:   fs.name,
		digest: d,
		outs:   []string{fs.out},
	}, nil
}
