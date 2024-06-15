package main

import (
	"github.com/podinate/podinate/cmd"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

func main() {
	cmd.Execute()
}
