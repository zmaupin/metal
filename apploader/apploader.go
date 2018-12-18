package apploader

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/metal-go/metal/apploader/static"
	"github.com/metal-go/metal/config"
)

// Run runs the App
func Run() error {
	// Channel of os.Signals for handling termination and channel for notification
	// event that we are done
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// Flush bundle from variable to disk. If any errors, send  interrupt signal
	// to trigger cleanup event and return error
	dir, err := flushBundle()

	log.Debugf("static assets under %s", dir)

	go func() {
		<-sigs
		os.RemoveAll(dir)
		done <- true
	}()
	if err != nil {
		sigs <- os.Interrupt
		<-done
		return err
	}

	// Create mux and configure handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Debug(fmt.Sprintf("serving request %v", r))
		handler := http.FileServer(http.Dir(path.Join(dir, "app", "dist")))
		handler.ServeHTTP(w, r)
	})

	// create run func that runs ListenAndServe in the background, sending errors
	// on the channel returned from it
	network := fmt.Sprintf("%s:%s", config.AppGlobal.GetAddress(), config.AppGlobal.GetPort())
	run := func() chan error {
		ch := make(chan error)
		go func() {
			ch <- http.ListenAndServe(network, mux)
		}()
		return ch
	}

	// Watch for SIGINT and SIGTERM signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// If an error occurs, return it. If we're done, return nil
	select {
	case err = <-run():
		return err
	case <-done:
		return nil
	}
}

// Crawl the tar archive of static assets and write files to disk
func flushBundle() (string, error) {
	tmpDir, err := ioutil.TempDir("", "metal-app")
	if err != nil {
		return tmpDir, err
	}
	bundle, err := base64.StdEncoding.DecodeString(strings.Trim(static.Bundle, "\""))
	static.Bundle = ""
	if err != nil {
		return tmpDir, err
	}
	gzr, err := gzip.NewReader(bytes.NewBuffer([]byte(bundle)))
	if err != nil {
		return tmpDir, err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return tmpDir, nil

		case err != nil:
			return tmpDir, err

		case header == nil:
			continue
		}

		target := filepath.Join(tmpDir, header.Name)
		log.Debug(fmt.Sprintf("unarchiving %s", target))

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return tmpDir, err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return tmpDir, err
			}
			if _, err := io.Copy(f, tr); err != nil {
				return tmpDir, err
			}
			f.Close()
		}
	}
}
