package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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
	logger *log.Logger
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

	startWebServerMux()
}

func startWebServerMux() {
	router := mux.NewRouter()
	router.Use(otelmux.Middleware(serviceName))

	cnt := 0

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("received request: #%d", cnt)

		cnt++

		// if cnt > 5 {
		// 	panic("Waaaaaaaaa !!!")
		// }

		body := fmt.Sprintf("Hello, world! %d", cnt)
		w.Write([]byte(body))
	})

	// router.HandleFunc("/healthcheck/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// })

	logger.Infoln("Running server on", endpoint)
	http.Handle("/", router)
	err := http.ListenAndServe(endpoint, nil)
	if err != nil {
		logger.Fatal(err)
	}
}
