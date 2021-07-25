package elsabin

import (
	"shanhu.io/elsa"
	"shanhu.io/misc/errcode"
)

const buildFile = "build.jsonx"
const buildSumsFile = "sums.jsonx"

func cmdSync(args []string) error {
	flags := cmdFlags.New()
	config := new(elsa.Config)
	declareBuildFlags(flags, config)
	pull := flags.Bool("pull", false, "pull latest commit")
	flags.ParseArgs(args)

	b := elsa.NewBuilder(config)

	build, err := elsa.ReadBuild(buildFile)
	if err != nil {
		return errcode.Annotate(err, "read build")
	}
	var sums *elsa.BuildSums
	if !*pull {
		s, err := elsa.ReadBuildSums(buildSumsFile)
		if err != nil {
			return errcode.Annotate(err, "read build sums")
		}
		sums = s
	}

	return b.SyncRepos(build, sums)
}
