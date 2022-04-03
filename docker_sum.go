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
	"strings"

	"shanhu.io/virgo/dock"
)

type dockerSum struct {
	ID     string
	Digest string
}

func newDockerSum(info *dock.ImageInfo, repo, prefer string) *dockerSum {
	sum := &dockerSum{ID: info.ID}

	var digests []string
	foundPrefered := false
	if repo != "" {
		digestPrefix := repo + "@"
		for _, d := range info.RepoDigests {
			if strings.HasPrefix(d, digestPrefix) {
				d = strings.TrimPrefix(d, digestPrefix)
				if d == prefer {
					foundPrefered = true
					break
				}
				digests = append(digests, d)
			}
		}
	}

	if foundPrefered {
		sum.Digest = prefer
	} else if len(digests) > 0 {
		sum.Digest = digests[0]
	}

	return sum
}
