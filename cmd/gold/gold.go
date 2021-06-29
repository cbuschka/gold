package main

import (
	"github.com/cbuschka/gold/internal/daemon"
	"github.com/kataras/golog"
	"os"
)

func main() {
	err := daemon.Run("./gold.conf.json")
	if err != nil {
		golog.Fatalf("Fatal error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
