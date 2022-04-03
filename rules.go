package caco3

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
