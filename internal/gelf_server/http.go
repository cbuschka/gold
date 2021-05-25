package gelf_server

import (
	jsonPkg "encoding/json"
	"fmt"
	journalPkg "github.com/cbuschka/golf/internal/journal"
	"github.com/gorilla/mux"
	gelf "gopkg.in/Graylog2/go-gelf.v2/gelf"
	"net"
	"net/http"
)

func ServeHttp(addr string, journal *journalPkg.Journal) error {

	httpListener, err := net.Listen("tcp", addr)
	fmt.Printf("HTTP server listening on %s...\n", addr)
	if err != nil {
		return err
	}
	defer httpListener.Close()
	httpHandler := newHttpHandler(journal)
	http.Serve(httpListener, httpHandler)

	return nil
}

func newHttpHandler(journal *journalPkg.Journal) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/gelf", func(w http.ResponseWriter, r *http.Request) {

		var gelfMessage gelf.Message
		err := jsonPkg.NewDecoder(r.Body).Decode(&gelfMessage)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		senderHost := r.Header.Get("X-Forwarded-For")
		if senderHost == "" {
			senderHost = r.RemoteAddr
		}
		message := journalPkg.FromGelfMessage(&gelfMessage, senderHost, "http")

		err = journal.WriteMessage(message)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Error(w, "", http.StatusCreated)
	}).Methods("POST")

	return http.Handler(router)
}
