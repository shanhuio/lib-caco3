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

const (
	ruleFileSet = "file_set"
	ruleBundle  = "bundle"
)

// FileSet selects a set of files.
type FileSet struct {
	Name string

	// The list of files to include in the fileset.
	Files []string `json:",omitempty"`

	// Selects a set of source input files.
	Select []string `json:",omitempty"`

	// Ignores a set of source input files after selection.
	Ignore []string `json:",omitempty"`

	// Merge in other file sets
	Include []string `json:",omitempty"`
}

// Bundle is a set of build rules in a bundle. A bundle has no build actions;
// it just group rules together.
type Bundle struct {
	// Name of the fule
	Name string

	// Other rule names.
	Deps []string
}
