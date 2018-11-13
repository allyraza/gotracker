package hyena

import (
	"net/http"

	// "github.com/Shopify/sarama"
	"github.com/cactus/go-statsd-client/statsd"
)

// Collector - state of the collector
type Collector struct {
	stats  statsd.Statter
	Buffer *buffer
	config *Config
	mux    *http.ServeMux
}

// NewCollector - creates a new collector
func NewCollector(config *Config) *Collector {
	// stats, err := statsd.NewClient(config.Stats.Addr, config.Stats.Key)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stats.Close()

	c := &Collector{}
	c.Buffer = NewBuffer(config)
	// c.stats = stats

	c.Routes()

	return c
}

// ServeHTTP - implements handler interface for directly using collector with http server
func (c *Collector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mux.ServeHTTP(w, r)
}
