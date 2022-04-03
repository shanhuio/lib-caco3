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
	"os"
)

type fileStat struct {
	Name         string
	Type         string
	Size         int64
	ModTimestamp int64
	Mode         uint32
}

const (
	fileTypeSrc = "s"
	fileTypeOut = "o"
)

func newOutFileStat(env *env, p string) (*fileStat, error) {
	return newFileStat(env, p, fileTypeOut)
}

func newSrcFileStat(env *env, p string) (*fileStat, error) {
	return newFileStat(env, p, fileTypeSrc)
}

func newFileStat(env *env, p, t string) (*fileStat, error) {
	var f string
	if t == fileTypeOut {
		f = env.out(p)
	} else {
		f = env.src(p)
	}

	info, err := os.Lstat(f)
	if err != nil {
		return nil, err
	}

	return &fileStat{
		Name:         p,
		Type:         t,
		Size:         info.Size(),
		ModTimestamp: info.ModTime().UnixNano(),
		Mode:         uint32(info.Mode()),
	}, nil
}
