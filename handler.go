package hyena

import (
	"log"
	"net/http"
)

// HandleClick it handles clicks
func (c *Collector) HandleClick(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	event := &Event{}
	event.Parse(r.URL.Query())

	log.Println(r.URL.String())

	go func() {
		c.Buffer.Put(event)
	}()

	w.WriteHeader(http.StatusOK)
}
