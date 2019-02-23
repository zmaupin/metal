package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/metal-go/metal/config"
)

func newRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleFunc("/streamoutput", serveStreamOutputClient)
	return router
}

func run() error {
	s := &http.Server{
		Addr:         config.RexecdGlobal.APIAddress,
		Handler:      newRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return s.ListenAndServe()
}
