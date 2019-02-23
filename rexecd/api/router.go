package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/metal-go/metal/config"
)

func newRouter(timeout time.Duration) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleFunc("/streamoutput", func(w http.ResponseWriter, r *http.Request) { serveStreamOutputClient(timeout, w, r) })
	return router
}

func run(timeout time.Duration) error {
	s := &http.Server{
		Addr:         config.RexecdGlobal.APIAddress,
		Handler:      newRouter(timeout),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}
	return s.ListenAndServe()
}
