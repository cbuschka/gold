package command_server

import (
	journalPkg "github.com/cbuschka/golf/internal/journal"
	"github.com/kataras/golog"
	"net"
	"net/http"
	"os"
)

func ServeUds(file string, journal journalPkg.Journal) error {

	if err := os.RemoveAll(file); err != nil {
		return err
	}

	golog.Infof("Command server http listener listening on %s...", file)

	udsListener, err := net.Listen("unix", file)
	if err != nil {
		return err
	}
	defer udsListener.Close()
	httpHandler := newHttpHandler(journal)
	err = http.Serve(udsListener, httpHandler)
	return err
}
