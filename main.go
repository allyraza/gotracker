package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "strings"

    "github.com/Shopify/sarama"
    "github.com/cactus/go-statsd-client/statsd"
)

// Config : tracker config
type Config struct {
    Addr         string
    StatsdAddr   string
    KafkaBrokers string
    Verbose      bool
}

// Server : state of the server
type Server struct {
    producer sarama.SyncProducer
    statsd   statsd.Statter
    config   *Config
}

// Run run server
func (s *Server) Run() {
    brokerList := strings.Split(s.config.KafkaBrokers, ",")

    // For the data collector, we are looking for strong consistency semantics.
    // Because we don't change the flush settings, sarama will try to produce messages
    // as fast as possible to keep latency low.
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
    config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
    config.Producer.Return.Successes = true

    // On the broker side, you may want to change the following settings to get
    // stronger consistency guarantees:
    // - For your broker, set `unclean.leader.election.enable` to false
    // - For the topic, you could increase `min.insync.replicas`.

    producer, err := sarama.NewSyncProducer(brokerList, config)
    if err != nil {
        log.Fatalln("Failed to start Sarama producer:", err)
    }
    s.producer = producer

    // The basic client sends one stat per packet (for compatibility).
    statsd, err := statsd.NewClient(s.config.StatsdAddr, "kafka")
    if err != nil {
        log.Fatal(err)
    }
    defer statsd.Close()
    s.statsd = statsd

    http.HandleFunc("/", s.HandleClick)
    http.ListenAndServe(s.config.Addr, nil)
}

// HandleClick it handles clicks
func (s *Server) HandleClick(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    // We are not setting a message key, which means that all messages will
    // be distributed randomly over the different partitions.
    partition, offset, err := s.producer.SendMessage(&sarama.ProducerMessage{
        Topic: "transactions",
        Value: sarama.StringEncoder(r.URL.RawQuery),
    })

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        if s.config.Verbose {
            fmt.Fprintf(w, "Failed to store your data:, %s\n", err)
            fmt.Printf("Failed to store your data:, %s\n", err)
        }
    } else {
        if s.config.Verbose {
            fmt.Fprintf(w, "Your data is stored with unique identifier: partition=%d, offset=%d\n", partition, offset)
            fmt.Printf("Your data is stored with unique identifier: partition=%d, offset=%d\n", partition, offset)
        }
        s.statsd.Inc("transaction_count", 1, 1)
    }
}

func main() {
    config := &Config{}

    flag.StringVar(&config.Addr, "addr", ":8080", "The address to bind to")
    flag.StringVar(&config.StatsdAddr, "statsd", "192.168.60.106:8125", "The statsd server address to connect to")
    flag.StringVar(&config.KafkaBrokers, "brokers", "0.0.0.0:9092", "The Kafka brokers to connect to, as a comma separated list")
    flag.BoolVar(&config.Verbose, "verbose", false, "Print out logging information for debugging purposes.")
    flag.Parse()

    server := &Server{config: config}
    server.Run()
}
