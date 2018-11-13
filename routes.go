package hyena

import (
	"net/http"
)

func (c *Collector) Routes() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", c.HandleClick)

	c.mux = mux
}
