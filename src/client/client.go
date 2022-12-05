package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/unvelista/loggingPlayground/src/logging"
)

const (
	defaultHostname = "0.0.0.0"
	defaultPort     = "8000"
	reqInterval     = 1 // seconds
)

var (
	logger *logrus.Logger
)

func main() {
	logger = logging.InitLogger()

	hostname, ok := os.LookupEnv("SERVER_HOSTNAME")
	if !ok {
		hostname = defaultHostname
	}
	port, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		port = defaultPort
	}
	url := fmt.Sprintf("http://%s:%s/", hostname, port)

	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		logger.Errorf("error building request: %w", err)
		syscall.Exit(1)
	}

	client := http.Client{}

	for {
		runRequests(req, &client)
		time.Sleep(reqInterval * time.Second)
	}
}

func runRequests(req *http.Request, client *http.Client) {
	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("error in request: %w", err)
		panic(err)
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Errorf(err.Error())
	}
	body := string(b)

	logger.Infoln(body)
}
