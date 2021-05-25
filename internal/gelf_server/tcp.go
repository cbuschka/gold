package gelf_server

import (
	"bufio"
	jsonPkg "encoding/json"
	"fmt"
	journalPkg "github.com/cbuschka/golf/internal/journal"
	worker "github.com/cbuschka/golf/internal/worker"
	gelf "gopkg.in/Graylog2/go-gelf.v2/gelf"
	"io"
	"net"
)

func ServeTcp(addr string, journal *journalPkg.Journal, workerPool *worker.WorkerPool) error {

	tcpListener, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	defer tcpListener.Close()

	fmt.Printf("Listening on %s/tcp...\n", addr)

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			return err
		}
		workerPool.RunWork(func() error {
			err := handleConnection(conn, journal)
			if err == io.EOF {
				return nil
			}
			return err
		})
	}

	return nil
}

func handleConnection(conn net.Conn, journal *journalPkg.Journal) error {
	fmt.Printf("Serving %s...\n", conn.RemoteAddr().String())
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		bbuf, err := readUntilZero(reader, 8192)
		if err != nil {
			return err
		}

		var gelfMessage gelf.Message
		err = jsonPkg.Unmarshal(bbuf, &gelfMessage)
		if err != nil {
			return err
		}

		message := journalPkg.FromGelfMessage(&gelfMessage, conn.RemoteAddr().String(), "tcp")
		err = journal.WriteMessage(message)
		if err != nil {
			return err
		}
	}

	return nil
}

func readUntilZero(reader *bufio.Reader, limit int) ([]byte, error) {
	bbuf := make([]byte, limit)
	var count = 0
	for {
		if count >= limit {
			return nil, fmt.Errorf("Limit %d exceeded.", limit)
		}

		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		if b == 0 {
			break
		}

		bbuf[count] = b
		count = count + 1
	}

	data := make([]byte, count)
	copy(data[:], bbuf[0:count])
	return data, nil
}
