package elsa

import (
	"shanhu.io/misc/jsonx"
)

// Build is the structure of the build.jsonx file. It specifies how
// to build a project.
type Build struct {
	Repos map[string]string
	Steps []*BuildStep `json:",omitempty"`
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
func ReadBuild(f string) (*Build, error) {
	b := new(Build)
	if err := jsonx.ReadFile(f, b); err != nil {
		return nil, err
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
