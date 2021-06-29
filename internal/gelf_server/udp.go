package gelf_server

import (
	journalPkg "github.com/cbuschka/gold/internal/journal"
	"github.com/kataras/golog"
	gelf "gopkg.in/Graylog2/go-gelf.v2/gelf"
)

func ServeUdp(addr string, journal journalPkg.Journal) error {
	rd, err := gelf.NewReader(addr)
	if err != nil {
		return err
	}

	golog.Infof("GELF udp listener listening on %s/udp...", addr)
	for {
		gelfMessage, err := rd.ReadMessage()
		if err != nil {
			return err
		}

		if gelfMessage == nil {
			break
		}

		message := journalPkg.FromGelfMessage(gelfMessage, "", "udp")
		err = journal.WriteMessage(message)
		if err != nil {
			golog.Errorf("Writing message %v failed: %v", message, err)
		}
	}

	return nil
}
