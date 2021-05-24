package main

import (
	"fmt"
	"github.com/cbuschka/golf/internal/daemon"
	"os"
)

func main() {
	err := daemon.Run()
	if err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
