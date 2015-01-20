package logbuddy

import (
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/mux"
)

//ServeStatic Serves static content
func (ws *WebServer) ServeStatic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["file"]
	data, err := Asset(strings.Join([]string{"static/", fileName}, ""))
	if err != nil {
		// Asset was not found.
		http.NotFound(w, r)
	}
	_, file := path.Split(fileName)
	fType := strings.Split(file, ".")
	if len(fType) == 2 {
		if fType[1] == "js" {
			w.Header().Set("Content-Type", "text/javascript")
		} else if fType[1] == "css" {
			w.Header().Set("Content-Type", "text/css")
		}
	}
	if len(fType) == 3 {
		if fType[2] == "js" {
			w.Header().Set("Content-Type", "text/javascript")
		} else if fType[2] == "css" {
			w.Header().Set("Content-Type", "text/css")
		}
	}
	w.Write(data)
}
