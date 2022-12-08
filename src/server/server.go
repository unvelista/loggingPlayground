package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/unvelista/loggingPlayground/src/logging"
	"github.com/unvelista/loggingPlayground/src/tracing"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

const (
	serviceName           = "server"
	endpoint              = "0.0.0.0:8000"
	enableTracing         = true
	otelCollectorEndpoint = "otel-collector:4317"
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
			logrus.Fatal(err)
		}
		defer tracing.FlushAndShutdownTracerProvider(context.Background(), tp)
	}

	startWebServerMux()
}

func startWebServerMux() {
	router := mux.NewRouter()
	router.Use(otelmux.Middleware(serviceName))

	cnt := 0

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.WithFields(logrus.Fields{
			"customField1": "bar",
			"customField2": 100,
		}).Infof("received request: #%d", cnt)

		cnt++

		// NOTE: for testing multiline parsing of stack traces
		// if cnt > 5 {
		// 	panic("Waaaaaaaaa !!!")
		// }

		body := fmt.Sprintf("Hello, world! %d", cnt)
		w.Write([]byte(body))
	})

	logger.Infoln("Running server on", endpoint)
	http.Handle("/", router)
	err := http.ListenAndServe(endpoint, nil)
	if err != nil {
		logger.Fatal(err)
	}
}
