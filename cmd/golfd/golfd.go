package main

import (
	"github.com/cbuschka/golf/internal/daemon"
	"github.com/kataras/golog"
	"os"
)

func main() {
	err := daemon.Run("./golfd.conf.json")
	if err != nil {
		golog.Fatalf("Fatal error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
