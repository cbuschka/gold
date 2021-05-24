package command_server

import (
	"fmt"
	journalPkg "github.com/cbuschka/golf/internal/journal"
	jsonPkg "encoding/json"
	gelf "gopkg.in/Graylog2/go-gelf.v2/gelf"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
)

func listMessages(w http.ResponseWriter, r *http.Request, begin int, limit int, journal *journalPkg.Journal) {
	fmt.Fprintf(w, "{\"messages\":[")
	isFirst := true
	journal.ListMessages(begin, limit, func(id uint64, message *gelf.Message) (bool, error) {
		messageWithId, err := toMessageWithId(id, message)
		if err != nil {
			return false, err
		}

		messageJson, err := jsonPkg.Marshal(messageWithId)
		if err != nil {
			return false, err
		}
		seperator := ",\n"
		if isFirst {
			seperator = "\n"
		}
		fmt.Fprintf(w, "%s%s", seperator, messageJson)
		isFirst = false
		return true, nil
	})
	fmt.Fprintf(w, "\n]}\n")
}

func newHttpHandler(journal *journalPkg.Journal) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		var begin = -1
		beginParam, ok := r.URL.Query()["begin"]
		if ok && beginParam[0] != "" {
			value, err := strconv.Atoi(beginParam[0])
			if  err != nil {
				http.Error(w, fmt.Sprintf("Begin invalid: '%s'", beginParam[0]), http.StatusBadRequest)
				return
			}
			begin = value
		}

		limit := -1
		limitParam, ok := r.URL.Query()["limit"]
		if ok && limitParam[0] != "" {
			value, err := strconv.Atoi(limitParam[0])
			if  err != nil {
				http.Error(w, fmt.Sprintf("Limit invalid: '%s'", limitParam[0]), http.StatusBadRequest)
				return
			}
			limit = value
		}

		listMessages(w, r, begin, limit, journal)
	}).Methods("GET")

	return http.Handler(router)
}

func toMessageWithId(id uint64, gelfMessage *gelf.Message) (*MessageWithId, error) {
	gelfMessageJson, err := jsonPkg.Marshal(gelfMessage)
	if err != nil {
		return nil, err
	}
	var messageWithId MessageWithId
	err = jsonPkg.Unmarshal(gelfMessageJson, &messageWithId)
	if err != nil {
		return nil, err
	}
	messageWithId.Id = id
	return &messageWithId, nil
}
