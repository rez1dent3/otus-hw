package main

import (
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(-1)
	}

	osEnv := os.Environ()
	env, err := ReadDir(os.Args[1])
	if err != nil {
		os.Exit(-1)
	}

	for _, str := range osEnv {
		values := strings.SplitN(str, "=", 2)
		if val, ok := env[values[0]]; ok {
			if val.NeedRemove {
				delete(env, str)
			}

			continue
		}

		env[values[0]] = EnvValue{Value: values[1], NeedRemove: false}
	}

	os.Exit(RunCmd(os.Args[2:], env))
}
