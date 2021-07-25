package elsabin

import (
	"shanhu.io/elsa"
	"shanhu.io/misc/flagutil"
)

var cmdFlags = flagutil.NewFactory("elsa")

func declareBuildFlags(flags *flagutil.FlagSet, c *elsa.Config) {
	flags.StringVar(&c.Src, "src", "src", "source directory")
	flags.StringVar(&c.Out, "out", "out", "output directory")
	flags.StringVar(&c.GoSrc, "gosrc", "", "go language source directory")
	flags.StringVar(&c.DockerRegistry, "cr", "cr.shanhu.io", "docker registry")
	flags.StringVar(
		&c.GoVersion, "goversion", "base/go-src.jsonx",
		"go language version spec file",
	)
	flags.StringVar(
		&c.SSHKnownHosts, "ssh_known_hosts", "base/ssh_known_hosts",
		"ssh known hosts file",
	)
}
