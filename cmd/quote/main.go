package main

import (
	"github.com/gorilla/mux"
	"github.com/mszostok/k8s-lecture/cmd/quote/internal/generator"
	"github.com/mszostok/k8s-lecture/cmd/quote/internal/web"
	"github.com/mszostok/k8s-lecture/internal/httperr"
	"github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Config represents mapping of ENV-s to in-app config structure
type Config struct {
	// Logger holds configuration for logger

	HTTPPort int `envconfig:"default=8080"`
	Quotes []string `envconfig:"optional"`
}

func main() {
	var cfg Config
	if err := envconfig.InitWithPrefix(&cfg, "APP"); err != nil {
		panic(err)
	}

	log := logrus.New().WithField("service", "quote")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/status", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})
	quoteProvider := &generator.QuoteProvider{
		Quotes: cfg.Quotes,
	}
	handlers := web.NewHandler(httperr.NewLogrusErrorReporter(log), quoteProvider)
	router.HandleFunc("/quote", handlers.GetRandomQuoteHandler)
	web.AddAPIRoutes(router, handlers)

	httpServer := http.Server{Addr: ":8080", Handler: router}

	go func() {
		log.Fatal(httpServer.ListenAndServe())
	}()

	log.Infof("Server started on port %d", cfg.HTTPPort)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}
