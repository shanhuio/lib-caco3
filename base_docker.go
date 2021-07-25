package elsa

type baseDocker struct {
	name  string
	build func(env *env, name string) error
}
