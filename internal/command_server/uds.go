package command_server

import (
	"fmt"
	journalPkg "github.com/cbuschka/golf/internal/journal"
	"net"
	"net/http"
	"os"
)

func ServeUds(file string, journal *journalPkg.Journal) {


	if err := os.RemoveAll(file); err != nil {
		panic(err)
	}

	fmt.Printf("Command server listening on %s...\n", file)

	udsListener, err := net.Listen("unix", file)
	if err != nil {
		panic(err)
	}
	defer udsListener.Close()
	httpHandler := newHttpHandler(journal)
	http.Serve(udsListener, httpHandler)
}
