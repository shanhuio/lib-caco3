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
	"shanhu.io/misc/jsonx"
	"shanhu.io/text/lexing"
)

// Workspace is the structure of the build.jsonx file. It specifies how
// to build a project.
type Workspace struct {
	RepoMap        *RepoMap
	Steps          []*BuildStep `json:",omitempty"`
	DockerSaveName bool         `json:",omitempty"`
}

// RepoMap contains the list of repos to clone down.
type RepoMap struct {
	GitHosting string `json:",omitempty"`
	Map        map[string]string
}

// BuildOptions contains the options to for the entire build.
type BuildOptions struct {
	DockerSaveName bool
}

// BuildStep is a rule for a step to build one or several targets in a
// directory.
type BuildStep struct {
	Name       string
	Dir        string       `json:",omitempty"`
	GoBinary   []string     `json:",omitempty"`
	NodeJS     *NodeJS      `json:",omitempty"`
	Dockers    []string     `json:",omitempty"`
	DockerPull *DockerPulls `json:",omitempty"`
}

// NodeJS is a rule to build a nodejs/npm package.
type NodeJS struct {
	Output []string `json:",omitempty"`
}

// DockerPulls specifies how to pull docker images from docker hub or other
// docker registries.
type DockerPulls struct {
	Images string
	Sums   map[string]string `json:",omitempty"`
}

// ReadWorkspace reads in the workspace manifest.
func ReadWorkspace(f string) (*Workspace, []*lexing.Error) {
	tm := func(t string) interface{} {
		switch t {
		case "repo_map":
			return new(RepoMap)
		case "build_step":
			return new(BuildStep)
		case "build_options":
			return new(BuildOptions)
		}
		return nil
	}
	entries, errs := jsonx.ReadSeriesFile(f, tm)
	if errs != nil {
		return nil, errs
	}

	ws := new(Workspace)
	for _, entry := range entries {
		switch v := entry.V.(type) {
		case *BuildStep:
			ws.Steps = append(ws.Steps, v)
		case *BuildOptions:
			ws.DockerSaveName = v.DockerSaveName
		case *RepoMap:
			ws.RepoMap = v
		}
	}
	return ws, nil
}

// RepoSums records the checkums and git commits of a build.
type RepoSums struct {
	RepoCommits map[string]string
}

// ReadRepoSums reads in the workspaces's repo checksum file.
func ReadRepoSums(f string) (*RepoSums, error) {
	b := new(RepoSums)
	if err := jsonx.ReadFile(f, b); err != nil {
		return nil, err
	}
	return b, nil
}

// SaveRepoSums saves sums to f.
func SaveRepoSums(f string, sums *RepoSums) error {
	return jsonx.WriteFile(f, sums)
}
