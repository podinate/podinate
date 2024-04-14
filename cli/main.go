package main

import (
	"github.com/johncave/podinate/cli/cmd"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

func main() {
	cmd.Execute()
}
