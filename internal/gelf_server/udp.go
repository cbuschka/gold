package gelf_server

import (
	"fmt"
	gelf "gopkg.in/Graylog2/go-gelf.v2/gelf"
	journalPkg "github.com/cbuschka/golf/internal/journal"
)

func ServeUdp(addr string, journal *journalPkg.Journal) error {
	rd, err := gelf.NewReader(addr)
	if err != nil {
		return err
	}

	fmt.Printf("Listening on %s/udp...\n", addr)
	for {
		message, err := rd.ReadMessage()
		if err != nil {
			return err
		}

		if message == nil {
			break;
		}

		err = journal.WriteMessage(message)
		if err != nil {
			fmt.Printf("Writing message %v failed.\n", message)
		} else {
			fmt.Printf("Message %v written to journal.\n", message)
		}
	}

	return nil
}
