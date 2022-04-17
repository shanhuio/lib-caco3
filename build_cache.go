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
	"encoding/json"
	"time"

	"shanhu.io/misc/errcode"
)

type buildCache struct {
}

type buildOutput struct {
	Src    *fileStat   `json:",omitempty"` // single source file.
	Outs   []*fileStat `json:",omitempty"` // Output S
	Docker *dockerSum  `json:",omitempty"`
}

func (c *buildCache) put(in string, out *buildOutput, t time.Time) error {
	bs, err := json.Marshal(out)
	if err != nil {
		return errcode.Annotate(err, "marshal output")
	}

	_ = bs
	panic("todo")
}
