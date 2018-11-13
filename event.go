package hyena

import (
	"net/url"
)

type Event struct {
	ID string `json:"id"`
	Type string `json:"type"`
}


func (e *Event) Parse(values url.Values) {
	e.ID = values.Get("id")
	e.Type = values.Get("t")
}