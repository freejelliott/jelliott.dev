package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ErrorHandler struct {
	fn     func(w http.ResponseWriter, r *http.Request) error
	logger *log.Logger
}

func (e ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := e.fn(w, r)
	if err != nil {
		e.logger.Printf("error serving request: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	}
}

func main() {
	logger := log.New(os.Stdout, "", log.Lshortfile)
	http.Handle("/", ErrorHandler{fn: home, logger: logger})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	port := 8080
	server := http.Server{
		Addr:     fmt.Sprintf(":%d", port),
		ErrorLog: logger,
	}
	logger.Printf("serving on port %d", port)
	err := server.ListenAndServe()
	logger.Fatalf("unexpected server shutdown: %s", err)
}

func home(w http.ResponseWriter, r *http.Request) error {
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(content))
	return err
}
