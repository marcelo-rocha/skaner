package main

import (
	"github.com/marcelo-rocha/skaner/cmd/skaner/cmd"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	cmd.Execute(logger)
}
