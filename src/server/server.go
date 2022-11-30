package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/unvelista/loggingPlayground/src/logging"
)

const (
	defaultHostname = "server"
	defaultPort     = "8000"
)

var (
	logger *log.Logger
)

func main() {
	fLogPath := "/server.log"
	f, err := os.OpenFile(fLogPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Failed to create logfile at %s: %+v", fLogPath, err)
		panic(err)
	}
	defer f.Close()

	logger = logging.InitLogger(f)

	startWebServerMux()
}

func startWebServerMux() {
	hostname, ok := os.LookupEnv("SERVER_HOSTNAME")
	if !ok {
		hostname = defaultHostname
	}
	port, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		port = defaultPort
	}
	addr := fmt.Sprintf("%s:%s", hostname, port)

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Infoln("received request")

		fmt.Fprintf(w, "Hello")
	})

	logger.Infoln("Running server on", addr)
	http.Handle("/", router)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Fatal(err)
	}
}
