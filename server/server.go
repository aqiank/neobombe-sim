package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackyb/neobombe-sim/simulation"
)

const Port = "8080"

type StateData struct {
	Units [simulation.NumUnits]simulation.Unit `json:"units"`
	Text  string                               `json:"text"`
}

var bombe simulation.Bombe
var text string

func indexHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	http.ServeFile(w, r, path)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	state := StateData{bombe.Units, text}
	if data, err := json.Marshal(state); err != nil {
		fmt.Println("server: updateHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(data)
	}
}

func Run(cs map[string]interface{}) {
	encChan := cs["encrypted"].(chan string)
	stateChan := cs["state"].(chan simulation.Bombe)
	sigChan := cs["signal"].(chan os.Signal)

	go func() {
		http.HandleFunc("/", indexHandler)
		http.HandleFunc("/update", updateHandler)
		fmt.Println("server starting..")
		log.Fatal(http.ListenAndServe(":"+Port, nil))
	}()
	go func() {
		for {
			text = <-encChan
		}
	}()
	go func() {
		for {
			bombe = <-stateChan
		}
	}()

	sigChan <- <-sigChan
	fmt.Println("server exiting..")
}
