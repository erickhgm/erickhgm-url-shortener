package config

import (
	"context"
	"sync"

	"ehgm.com.br/url-shortener/domain/ports"

	"cloud.google.com/go/firestore"
	"github.com/go-redis/redis/v8"
)

var (
	firestoreClient *firestore.Client
	redisClient     *redis.Client
)

func createFirestoreClient(ctx context.Context, log ports.Logger, projectId string) {
	var err error
	firestoreClient, err = firestore.NewClient(ctx, projectId)
	if err != nil {
		log.Fatal("Failed to create firestore client: %s", err)
	}
}

func createRedisClient(host, pass string) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pass,
		DB:       0,
	})
}

func NewFirestoreClient(ctx context.Context, log ports.Logger, projectId string) *firestore.Client {
	var once sync.Once
	once.Do(func() { createFirestoreClient(ctx, log, projectId) })
	return firestoreClient
}

func NewRedisClient(redisHost, redisPass string) *redis.Client {
	var once sync.Once
	once.Do(func() { createRedisClient(redisHost, redisPass) })
	return redisClient
}
