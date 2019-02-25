package api

import (
	"database/sql"
	"log"
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
func Run(done chan struct{}) error {
	// Get and set a sql.DB
	dsn := config.RexecdGlobal.DataSourceName + "rexecd"
	db, err := sql.Open("mysql", dsn)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:         config.RexecdGlobal.APIAddress,
		Handler:      newRouter(config.RexecdGlobal.APITimeout, db),
		ReadTimeout:  config.RexecdGlobal.APITimeout,
		WriteTimeout: config.RexecdGlobal.APITimeout,
	}
	run := func() chan error {
		errChan := make(chan error)
		go func() {
			errChan <- s.ListenAndServe()
		}()
		return errChan
	}
	select {
	case err := <-run():
		return err
	case <-done:
		return nil
	}
}
