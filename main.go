package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := &Config{}

	flag.StringVar(&config.Filepath, "config", "config.json", "The config file for collector.")
	flag.BoolVar(&config.Verbose, "verbose", false, "Print out debugging info.")
	flag.Parse()

	config.ParseFile()

	log.Printf("Starting server on %v", config.Addr)

	collector := NewCollector(config)

	server := http.Server{
		Addr:    config.Addr,
		Handler: collector,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	log.Println("Shutdown signal recieved, quiting...")
	server.Shutdown(context.Background())
}
