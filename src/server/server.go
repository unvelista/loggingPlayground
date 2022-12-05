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
	logger = logging.InitLogger()

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

	cnt := 0

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("received request: #%d", cnt)

		cnt++

		// if cnt > 5 {
		// 	panic("Waaaaaaaaa !!!")
		// }

		body := "Hello, world!"
		w.Write([]byte(body))
	})

	logger.Infoln("Running server on", addr)
	http.Handle("/", router)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Fatal(err)
	}
}
