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
