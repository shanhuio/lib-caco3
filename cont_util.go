package elsa

import (
	"fmt"
	"path"

	"shanhu.io/virgo/dock"
)

func exitError(exit int) error {
	if exit == 0 {
		return nil
	}
	return fmt.Errorf("exit with code: %d", exit)
}

func execError(ret int, err error) error {
	if err != nil {
		return err
	}
	return exitError(ret)
}

func contExec(cont *dock.Cont, args []string) error {
	return execError(cont.ExecArgs(args))
}

func linuxPathJoin(parts ...string) string {
	return path.Join(parts...)
}
