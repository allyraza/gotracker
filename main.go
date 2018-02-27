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

var (
	addr       = flag.String("addr", ":8080", "The address to bind to")
	statsdAddr = flag.String("statsd", "192.168.60.106:8125", "The statsd server address to connect to")
	brokers    = flag.String("brokers", "0.0.0.0:9092", "The Kafka brokers to connect to, as a comma separated list")
)

// HandleClick it handles clicks
func HandleClick(w http.ResponseWriter, r *http.Request) {
	brokerList := strings.Split(*brokers, ",")

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

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// We are not setting a message key, which means that all messages will
	// be distributed randomly over the different partitions.
	partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: "transactions",
		Value: sarama.StringEncoder(r.URL.RawQuery),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to store your data:, %s\n", err)
		fmt.Printf("Failed to store your data:, %s\n", err)
	} else {
		// The tuple (topic, partition, offset) can be used as a unique identifier
		// for a message in a Kafka cluster.
		fmt.Fprintf(w, "Your data is stored with unique identifier important/%d/%d\n", partition, offset)
		fmt.Printf("Your data is stored with unique identifier important/%d/%d\n", partition, offset)

	}
}

func incrementMetrics() {
	// first create a client
	// The basic client sends one stat per packet (for compatibility).
	client, err := statsd.NewClient(*statsdAddr, "kafka-client")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	client.Inc("kafka_count", 1, 1.0)
}

func main() {
	flag.Parse()
	http.HandleFunc("/", HandleClick)
	http.ListenAndServe(*addr, nil)
}
