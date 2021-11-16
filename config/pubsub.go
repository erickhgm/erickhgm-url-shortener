package config

import (
	"context"
	"sync"

	"ehgm.com.br/url-shortener/domain/ports"

	"cloud.google.com/go/pubsub"
)

var pubsubClient *pubsub.Client

func createPubSubClient(ctx context.Context, log ports.Logger, projectId string) {
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		log.Fatal("Failed to create pubsub client: %s", err)
	}
	pubsubClient = client
}

func NewPubSubClient(ctx context.Context, log ports.Logger, projectId string) *pubsub.Client {
	var once sync.Once
	once.Do(func() { createPubSubClient(ctx, log, projectId) })
	return pubsubClient
}
