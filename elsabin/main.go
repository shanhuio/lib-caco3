package elsabin

import (
	"shanhu.io/misc/subcmd"
)

func cmd() *subcmd.List {
	c := subcmd.New()
	c.Add("build", "builds a target", cmdBuild)
	c.Add("sync", "sync source repos", cmdSync)
	return c
}

// Main is the entrance for the elsa binary.
func Main() { cmd().Main() }
