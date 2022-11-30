package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"
)

const (
	defaultHostname = "0.0.0.0"
	defaultPort     = "8000"
	reqInterval     = 1 // seconds
)

func main() {
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
		fmt.Errorf("error building request: %w", err)
		syscall.Exit(1)
	}

	client := http.Client{}

	for {
		res, err := client.Do(req)
		if err != nil {
			fmt.Errorf("error in request: %w", err)
			panic(err)
		}
		log.Println(res.Body)
		time.Sleep(reqInterval * time.Second)
	}
}
