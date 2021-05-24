package command_server

import (
	"fmt"
	journalPkg "github.com/cbuschka/golf/internal/journal"
	"net"
	"net/http"
	"os"
)

func ServeUds(file string, journal *journalPkg.Journal) error {

	if err := os.RemoveAll(file); err != nil {
		return err
	}

	fmt.Printf("Command server listening on %s...\n", file)

	udsListener, err := net.Listen("unix", file)
	if err != nil {
		return err
	}
	defer udsListener.Close()
	httpHandler := newHttpHandler(journal)
	http.Serve(udsListener, httpHandler)

	return nil
}
