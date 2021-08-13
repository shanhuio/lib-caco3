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
	save := flags.Bool("save", false, "save latest commit into sums file")
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

	newSums, err := b.SyncRepos(build, sums)
	if err != nil {
		return err
	}
	if *save {
		if err := elsa.SaveBuildSums(buildSumsFile, newSums); err != nil {
			return errcode.Annotate(err, "save build sums")
		}
	}
	return nil
}
