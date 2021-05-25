package gelf_server

import (
	"fmt"
	journalPkg "github.com/cbuschka/golf/internal/journal"
	gelf "gopkg.in/Graylog2/go-gelf.v2/gelf"
)

func ServeUdp(addr string, journal *journalPkg.Journal) error {
	rd, err := gelf.NewReader(addr)
	if err != nil {
		return err
	}

	fmt.Printf("Listening on %s/udp...\n", addr)
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
			fmt.Printf("Writing message %v failed: %v\n", message, err)
		}
	}

	return nil
}
