package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/allyraza/hyena"
)

func main() {
	config := &hyena.Config{}

	flag.StringVar(&config.Filepath, "config", "config.json", "The config file for collector.")
	flag.BoolVar(&config.Verbose, "verbose", false, "Print out debugging info.")
	flag.Parse()

	config.ParseFile()

	log.Printf("Starting server on %v\n", config.Addr)

	collector := hyena.NewCollector(config)

	worker := hyena.NewWorker(collector.Buffer)
	go func() {
		worker.Start()
	}()

	server := http.Server{Addr: config.Addr, Handler: collector}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	log.Println("Shutdown signal received, exiting...")
	server.Shutdown(context.Background())
}
