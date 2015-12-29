package libs

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func ExecCmd(path string, arg1 string, args ...string) (string, error) {
	var cmd *exec.Cmd

	argsNew := make([]string, len(args))

	for i, a := range args {
		argsNew[i] = a
	}

	cmd = exec.Command(arg1, argsNew...)
	cmd.Dir = path
	// cmd.stdin = os.Stdin
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	var retStr string
	err = cmd.Start()
	if err != nil {
		retb, _ := ioutil.ReadAll(stderr)
		retStr = string(retb)
	} else {
		retb, _ := ioutil.ReadAll(stdout)
		retStr = string(retb)
	}

	return retStr, err
}
