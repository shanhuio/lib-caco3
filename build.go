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

package elsa

import (
	"shanhu.io/misc/jsonx"
	"shanhu.io/text/lexing"
)

// Build is the structure of the build.jsonx file. It specifies how
// to build a project.
type Build struct {
	Repos          map[string]string
	Steps          []*BuildStep `json:",omitempty"`
	DockerSaveName bool         `json:",omitempty"`
}

// RepoMap contains the list of repos to clone down.
type RepoMap struct {
	Map map[string]string
}

// BuildOptions contains the options to for the entire build.
type BuildOptions struct {
	DockerSaveName bool
}

// BuildStep is a rule for a step to build one or several targets in a
// directory.
type BuildStep struct {
	Name       string
	Dir        string      `json:",omitempty"`
	GoBinary   []string    `json:",omitempty"`
	NodeJS     *NodeJS     `json:",omitempty"`
	Dockers    []string    `json:",omitempty"`
	DockerPull *DockerPull `json:",omitempty"`
}

// NodeJS is a rule to build a nodejs/npm package.
type NodeJS struct {
	Output []string `json:",omitempty"`
}

// DockerPull specifies how to pull docker images from docker hub or other
// docker registries.
type DockerPull struct {
	Images string
	Sums   map[string]string `json:",omitempty"`
}

// ReadBuild reads in a build manifest.
func ReadBuild(f string) (*Build, []*lexing.Error) {
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

	b := new(Build)
	for _, entry := range entries {
		switch v := entry.V.(type) {
		case *BuildStep:
			b.Steps = append(b.Steps, v)
		case *BuildOptions:
			b.DockerSaveName = v.DockerSaveName
		case *RepoMap:
			b.Repos = v.Map
		}
	}
	return b, nil
}

// BuildSums records the checkums and git commits of a build.
type BuildSums struct {
	RepoCommits map[string]string
}

// ReadBuildSums reads in the build's checksum file.
func ReadBuildSums(f string) (*BuildSums, error) {
	b := new(BuildSums)
	if err := jsonx.ReadFile(f, b); err != nil {
		return nil, err
	}
	return b, nil
}

// SaveBuildSums saves sums to f.
func SaveBuildSums(f string, sums *BuildSums) error {
	return jsonx.WriteFile(f, sums)
}
