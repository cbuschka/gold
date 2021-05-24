package command_server

import (
	"fmt"
	journalPkg "github.com/cbuschka/golf/internal/journal"
	jsonPkg "encoding/json"
	gelf "gopkg.in/Graylog2/go-gelf.v2/gelf"
	"net/http"
	"github.com/gorilla/mux"
)

type listMessagesResponse struct {
	messages []*gelf.Message
}

func listMessages(w http.ResponseWriter, r *http.Request, journal *journalPkg.Journal) {
	fmt.Fprintf(w, "{\"messages\":[")
	first := true
	journal.ListMessages(func(id uint64, message *gelf.Message) (bool, error) {
		messageJson, err := jsonPkg.Marshal(message)
		if err != nil {
			return false, err
		}
		seperator := ",\n"
		if first {
			seperator = "\n"
		}
		fmt.Fprintf(w, "%s%s", seperator, messageJson)
		first = false
		return true, nil
	})
	fmt.Fprintf(w, "\n]}\n")
}

func newHttpHandler(journal *journalPkg.Journal) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		listMessages(w, r, journal)
	}).Methods("GET")

	return http.Handler(router)
}
