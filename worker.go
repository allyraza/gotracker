package hyena

import (
	"fmt"
	"log"
	"encoding/json"
)

type worker struct {
	config *Config
	buffer *buffer
}

func NewWorker(buffer *buffer) *worker {
	return &worker{
		buffer: buffer,
	}
}

func (c *worker) Start() {
	for {
		event, err := c.buffer.Get()
		if err != nil {
			log.Printf("Worker: %v\n", err)
		}

		blob, _ := json.Marshal(event)
		fmt.Printf("%v\n", string(blob))
	}

	// We are not setting a message key, which means that all messages will
	// be distributed randomly over the different partitions.
	// partition, offset, err := c.buffer.SendMessage(&sarama.ProducerMessage{
	// 	Topic: "transactions",
	// 	Value: sarama.StringEncoder(r.URL.RawQuery),
	// })

	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	if c.config.Verbose {
	// 		fmt.Fprintf(w, "Failed to store your data:, %s\n", err)
	// 		fmt.Printf("Failed to store your data:, %s\n", err)
	// 	}
	// } else {
	// 	if c.config.Verbose {
	// 		fmt.Fprintf(w, "Your data is stored with unique identifier: partition=%d, offset=%d\n", partition, offset)
	// 		fmt.Printf("Your data is stored with unique identifier: partition=%d, offset=%d\n", partition, offset)
	// 	}
	// 	c.stats.Inc("transaction_count", 1, 1)
	// }

}
