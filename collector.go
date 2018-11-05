package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/cactus/go-statsd-client/statsd"
)

// Collector - state of the collector
type Collector struct {
	stats  statsd.Statter
	buffer sarama.SyncProducer
	config *Config
	mux    *http.ServeMux
}

// Run run server
func NewCollector(config *Config) *Collector {

	fmt.Printf("%v\n", config.Kafka)

	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	sc := sarama.NewConfig()
	sc.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	sc.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	sc.Producer.Return.Successes = true

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewSyncProducer(config.Kafka.Brokers, sc)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	// The basic client sends one stat per packet (for compatibility).
	stats, err := statsd.NewClient(config.Stats, "kafka")
	if err != nil {
		log.Fatal(err)
	}
	defer stats.Close()

	c := &Collector{}
	c.buffer = producer
	c.stats = stats

	c.Routes()

	return c
}

func (c *Collector) Routes() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", c.HandleClick)

	c.mux = mux
}

func (c *Collector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mux.ServeHTTP(w, r)
}

// HandleClick it handles clicks
func (c *Collector) HandleClick(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// We are not setting a message key, which means that all messages will
	// be distributed randomly over the different partitions.
	partition, offset, err := c.buffer.SendMessage(&sarama.ProducerMessage{
		Topic: "transactions",
		Value: sarama.StringEncoder(r.URL.RawQuery),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if c.config.Verbose {
			fmt.Fprintf(w, "Failed to store your data:, %s\n", err)
			fmt.Printf("Failed to store your data:, %s\n", err)
		}
	} else {
		if c.config.Verbose {
			fmt.Fprintf(w, "Your data is stored with unique identifier: partition=%d, offset=%d\n", partition, offset)
			fmt.Printf("Your data is stored with unique identifier: partition=%d, offset=%d\n", partition, offset)
		}
		c.stats.Inc("transaction_count", 1, 1)
	}
}
