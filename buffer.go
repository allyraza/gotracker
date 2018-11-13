package hyena

import "errors"

type buffer struct {
	queue chan *Event
}

// sc := sarama.NewConfig()
// sc.Producer.RequiredAcks = sarama.WaitForAll
// sc.Producer.Retry.Max = 10
// sc.Producer.Return.Successes = true
//
// producer, err := sarama.NewSyncProducer(config.Kafka.Brokers, sc)
// if err != nil {
// 	log.Fatalln("Failed to start Sarama producer:", err)
// }

func NewBuffer(config *Config) *buffer {
	return &buffer{
		queue: make(chan *Event, config.Buffer.Size),
	}
}

func (b *buffer) Put(event *Event) error {
	if event == nil {
		return errors.New("Buffer: can't put invalid event")
	}

	b.queue <- event

	return nil
}

func (b *buffer) Get() (*Event, error) {
	select {
	case event := <-b.queue:
		return event, nil
	}

	return nil, errors.New("Buffer: no events available")
}
