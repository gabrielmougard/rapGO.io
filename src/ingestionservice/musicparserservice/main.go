package main

import (
	"net/http"
	"time"
	"fmt"

	"github.com/gorilla/mux"

	"rapGO.io/src/ingestionservice/musicparserservice/api"
	"rapGO.io/src/ingestionservice/musicparserservice/pkg/setting"


)

func main() {

	//Other interesting genre :
	// "Soundtrack", "Nerdcore", "Wonky"
	// baseURLs := []string{"Hip-Hop_Beats","Hip-Hop"}
	// var mp MusicParser
	// mp.Seed(baseURLs,true,23)
	// mp.Start()
	router := mux.NewRouter()
	router.HandleFunc("/ingest/url", api.handleURL).Methods("POST")
	router.HandleFunc("/ingest/urlfile", api.handleURLfile).Methods("POST")
	router.HandleFunc("/ingest/rawfile", api.handleRawFile).Methods("POST")
	router.HandleFunc("/ingest/genre/{genre}", api.handleGenreParsing).Methods("POST")

	router.HandleFunc("/ingest/genre", api.handleGetGenre).Methods("GET")

	//TODO : Websocket route(s) or GET route (with a poll delay specified in the front part)
	router.HandleFunc("/bucket/data", api.handleGetBucketData).Methods("GET")
	//

	fmt.Println("Starting server on port "+setting.getServerPort()+" ...")
	http.ListenAndServe(":"+setting.getServerPort(), nil)
}