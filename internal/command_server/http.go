package command_server

import (
	jsonPkg "encoding/json"
	"fmt"
	journalPkg "github.com/cbuschka/golf/internal/journal"
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"net/http"
	"strconv"
)

func listMessages(w http.ResponseWriter, r *http.Request, begin string, limit int, journal journalPkg.Journal) error {
	_, err := fmt.Fprintf(w, "{\"messages\":[")
	if err != nil {
		return err
	}
	isFirst := true
	err = journal.ListMessages(begin, limit, func(message *journalPkg.Message) (bool, error) {

		messageJson, err := jsonPkg.Marshal(message)
		if err != nil {
			return false, err
		}

		seperator := ",\n"
		if isFirst {
			seperator = "\n"
		}
		_, err = fmt.Fprintf(w, "%s%s", seperator, messageJson)
		if err != nil {
			return false, err
		}

		isFirst = false
		return true, nil
	})
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "\n]}\n")

	return err
}

func newHttpHandler(journal journalPkg.Journal) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		var begin = ""
		beginParam, ok := r.URL.Query()["begin"]
		if ok && len(beginParam) > 0 && beginParam[0] != "" {
			begin = beginParam[0]
		}

		limit := -1
		limitParam, ok := r.URL.Query()["limit"]
		if ok && len(limitParam) > 0 && limitParam[0] != "" {
			value, err := strconv.Atoi(limitParam[0])
			if err != nil {
				http.Error(w, fmt.Sprintf("Limit invalid: '%s'", limitParam[0]), http.StatusBadRequest)
				return
			}
			limit = value
		}

		err := listMessages(w, r, begin, limit, journal)
		if err != nil {
			golog.Errorf("Listing messages failed: %v", err)
		}
	}).Methods("GET")

	return http.Handler(router)
}
