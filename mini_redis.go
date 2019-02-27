package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	store := new(Store)

	go serveHttp(store)
	runShell(store)
}

const defaultHttpPort = "8080"

func serveHttp(store *Store) {
	handler := HttpHandler{Interpreter{store}}

	addr := fmt.Sprintf(":%s", defaultHttpPort)
	err := http.ListenAndServe(addr, handler)
	log.Fatal(err)
}

type HttpHandler struct {
	Interpreter
}

func (handler HttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var serveErr error

	if cmd := req.FormValue("cmd"); cmd == "" {
		serveErr = respondJson(w, http.StatusBadRequest, map[string]string{
			"error": "No valid \"cmd\" query parameter identified",
		})
	} else {
		if value, err := handler.Exec(cmd); err == nil {
			serveErr = respondJson(w, http.StatusOK, value)
		} else {
			serveErr = respondJson(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
	}

	if serveErr != nil {
		log.Fatal("Got the following error while serving http request: ", serveErr)
	}
}

func respondJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func runShell(store *Store) {
	intr := Interpreter{store}
	fmt.Println("Type \"exit\" to leave")

	scanner := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")

		txt, err := scanner.ReadString('\n')
		if err != nil {
			fmt.Println("Got the following error while retrieving input: ", err)
			continue
		}

		cmd := strings.TrimSpace(txt)

		switch cmd {
		case "":
		case "exit":
			fmt.Println("Exiting...")
			return
		default:
			if actual, err := intr.Exec(cmd); err == nil {
				fmt.Printf("  %v\n", actual)
			} else {
				fmt.Println("Command failed with the following error: ", err)
			}
		}
	}
}
