package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/unvelista/loggingPlayground/src/logging"
	"github.com/unvelista/loggingPlayground/src/tracing"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	serviceName           = "client"
	serverEndpoint        = "server:8000"
	enableTracing         = true
	otelCollectorEndpoint = "otel-collector:4317"
	reqInterval           = 1 // seconds

)

var (
	logger *logrus.Entry
)

func main() {
	logger = logging.InitLogger()

	if enableTracing {
		tp, err := tracing.InitTracerProvider(
			context.Background(), serviceName, otelCollectorEndpoint)
		if err != nil {
			log.Fatal(err)
		}
		defer tracing.FlushAndShutdownTracerProvider(context.Background(), tp)
	}

	logger.Infof("using %s as server endpoint", serverEndpoint)

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://%s/", serverEndpoint),
		nil,
	)
	if err != nil {
		logger.Errorf("error building request: %w", err)
		syscall.Exit(1)
	}

	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	for {
		runRequests(req, client)
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
