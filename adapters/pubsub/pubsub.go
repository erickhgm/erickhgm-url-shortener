package pubsub

import (
	"context"

	"ehgm.com.br/url-shortener/domain/ports"

	"cloud.google.com/go/pubsub"
)

// Struct that implements 'UrlCounter' interface
type urlCounter struct {
	log         ports.Logger
	ps          *pubsub.Client
	pubsubTopic string
}

// Get an instance of 'UrlCounter' using this method
func NewUrlCounter(log ports.Logger, ps *pubsub.Client, pubsubTopic string) ports.UrlCounter {
	return &urlCounter{log: log, ps: ps, pubsubTopic: pubsubTopic}
}

func (c *urlCounter) IncrementCounter(id string) {
	ctx := context.Background()
	topic := c.ps.Topic(c.pubsubTopic)
	result := topic.Publish(ctx, &pubsub.Message{Data: []byte(id)})

	idMessage, err := result.Get(ctx)
	if err != nil {
		c.log.Error("Error sending message to Id: %v. Cause: %s", id, err)
	} else {
		c.log.Info("Message [%v] sent successfully to Id: %v", idMessage, id)
	}
}
