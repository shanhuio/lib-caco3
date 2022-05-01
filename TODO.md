## TODO

- sum db and sum file update
- allow `BUILD.caco3` in sub repos and directories
- build all homedrv and shanhu artifacts
- add build cache

----

repo organization

````
shanhu.io
	base
		misc
		xxx
		xxx
	homedrv
		dockers
		homedrv
		jarvis-ts
		homedrv-ts
````

- container image: shanhu.io/proj/dockers -> cr.shanhu.io/proj/...
	- github: github.com/shanhuio/proj-dockers
- golang pkg: shanhu.io/proj/name -> shanhu.io/proj/name
	- github: github.com/shanhuio/proj-name
- npm package: shanhu.io/proj/name-ts -> @shanhu.io/proj-name
	- github: github.com/shanhuio/proj-name-ts
- other repos: shanhu.io/proj/name
	- github: github.com/shanhuio/proj-name

- this means shanhu.io/shanhu needs to go as shanhu.io/lab/shanhu
- and homedrv needs to go as
- which is a bit tedious, but is is 

````
current					changed				github				golang/npm
core/caco3-base			lib/dockers			lib-dockers
core/caco3				lib/caco3			lib-caco3			shanhu.io/caco3
core/misc				lib/misc			lib-misc            shanhu.io/misc
core/pisces				lib/pisces			lib-pisces          shanhu.io/pisces
core/aries				lib/aries			lib-aries           shanhu.io/aries
core/htmlgen-ts			lib/htmlgen-ts		lib-htmlgen-ts		@shanhuio/htmlgen
core/misc-ts			lib/misc-ts			lib-misc-ts			@shanhuio/misc
core/text				lib/text			lib-text
core/dags               lib/dags            lib-dags
core/dags-ts			lib/dags-ts			lib-dags-ts			@shanhuio/dags
core/lessbase			lib/style-ts		lib-style-ts		@shanhuio/style

homedrv/homedrv-ts 		homedrv/site-ts
homedrv/homedrv			homedrv/drv			homedrv-drv
homedrv/homedrv-build	homedrv/build		homedrv-build
homedrv/homedrv-dockers	homedrv/dockers		homedrv-dockers  // public dockers
homedrv/jarvis-ts		homedrv/drv-ts		homedrv-drv-ts
homedrv/tunlsite-ts		homedrv/tunl-ts
homedrv/homedrv-docs	homedrv/docs		homedrv-docs
						homedrv/core

smlrepo/smlrepo			smlrepo/core
smlrepo/sml				smlrepo/sml         
smlrepo/smlrepo-ts		smlrepo/site-ts     @shanhuio/smlrepo-site
smlrepo/dagvis-ts		smlrepo/dagvis-ts   @shanhuio/smlrepo-dagvis
smlrepo/tools			smlrepo/tools

smlvm/smlg				smlvm/site
smlvm/smlvm				smlvm/vm
smlvm/smlg-ts			smlvm/site-ts       @shanhuio/smlvm-site
smlvm/smldriod			smlvm/droid
smlvm/smlos				smlvm/os
smlvm/smlhome			smlvm/gbase
smlvm/gtutor			smlvm/gtutor

lab/tunlsite			homedrv/tunl
lab/shanhu				lab/shanhu			// need to split out homedrv stuff
lab/galaxy				lab/galaxy
lab/mp3tagfix			lab/mp3util
lab/gcimporter			lab/gcimporter
lab/build				lab/build			// dispatch to other projects
lab/shanhu-ts			corp/site-ts
lab/shanhu-dockers		lab/dockers			// dispatch to other projects
lab/tunlsite-dockers	homedrv/priv-dockers	// private dockers
lab/smlts				smlvm/?				// not sure if still in use
lab/docs				lab/docs
lab/third				lab/third
lab/yumuzi				lab/yumuzi

````

## migration plan

- find a tool to rename all imports
- get a list of all repositories
- halt development
- clone public repos
- rename private repos
- add new shanhu.io gometas
- rename imports of everything
	- build a script
- remove old repositories

----

# BUILD in repos

- We will scan the root of 