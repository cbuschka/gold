package main

import (
	"fmt"
	"github.com/cbuschka/golf/internal/query"
	"os"
)

func main() {
	err := query.ListMessages()
	if err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
