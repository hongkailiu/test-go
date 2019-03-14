package main

import (
	"fmt"
	"os"

	"github.com/hongkailiu/test-go/pkg/testctl/cmd"
)

func main() {
	command := cmd.NewDefaultTestctlCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
