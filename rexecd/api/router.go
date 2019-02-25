package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/metal-go/metal/config"
)

func newRouter(timeout time.Duration, db *sql.DB) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleFunc("/command", func(w http.ResponseWriter, r *http.Request) { command(db, w, r) })
	router.HandleFunc("/streamoutput", func(w http.ResponseWriter, r *http.Request) { serveStreamOutputClient(timeout, w, r) })
	return router
}

// Run is the runner for the Rexecd API server
func Run(timeout time.Duration, db *sql.DB) error {
	s := &http.Server{
		Addr:         config.RexecdGlobal.APIAddress,
		Handler:      newRouter(timeout, db),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}
	return s.ListenAndServe()
}
