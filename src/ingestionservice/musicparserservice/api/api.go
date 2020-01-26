package api

import (
	"fmt"
	"time"
	"net/http"

	"github.com/gorilla/mux"

	"rapGO.io/src/ingestionservice/musicparserservice/pkg/download"

)

func handleURL(w http.ResponseWriter, r *http.Request) {
	urls := r.
	download.downloadFile()
}

func handleURLfile(w http.ResponseWriter, r *http.Request) {
	
}

func handleRawFile(w http.ResponseWriter, r *http.Request) {
	
}

func handleGenreParsing(w http.ResponseWriter, r *http.Request) {
	
}