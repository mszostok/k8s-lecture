package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mszostok/k8s-lecture/cmd/meme/internal/generator"
	"github.com/mszostok/k8s-lecture/cmd/meme/internal/quote"
	"github.com/mszostok/k8s-lecture/cmd/meme/internal/web"
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
	HTTPPort        int    `envconfig:"PORT,default=9090"`
	QuoteServiceUrl string `envconfig:"QUOTE_URL"`
}

func main() {
	var cfg Config
	if err := envconfig.InitWithPrefix(&cfg, ""); err != nil {
		panic(err)
	}

	log := logrus.New().WithField("service", "meme")

	quoteClient := quote.NewClient(http.DefaultClient, cfg.QuoteServiceUrl)
	memGen := generator.New(quoteClient)
	h := web.NewHandler(memGen, httperr.NewLogrusErrorReporter(log))
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/status", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})
	web.AddAPIRoutes(router, h)

	httpServer := http.Server{Addr: fmt.Sprintf(":%d", cfg.HTTPPort), Handler: router}

	go func() {
		log.Fatal(httpServer.ListenAndServe())
	}()

	log.Info("Server started on port %d", cfg.HTTPPort)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}
