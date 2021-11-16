package config

import (
	"os"
	"strconv"

	"ehgm.com.br/url-shortener/domain/ports"
)

type EnvConfig struct {
	ProjectId   string
	RedisHost   string
	RedisPass   string
	RedisTTL    int
	PubsubTopic string
	IdLength    int
}

func NewEnvConfig(log ports.Logger) EnvConfig {
	log.Info("Starting NewEnvConfig ...")

	project := os.Getenv("PROJECT_ID")
	redisHost := os.Getenv("REDIS_HOST")
	redisPass := os.Getenv("REDIS_PASS")
	redisTTL := os.Getenv("REDIS_TTL")
	psTopic := os.Getenv("PUBSUB_TOPIC")
	idLenght := os.Getenv("ID_LENGHT")

	if len(project) <= 0 {
		log.Fatal("Failed to load PROJECT_ID environment variable")
	}
	if len(redisHost) <= 0 {
		log.Fatal("Failed to load REDIS_HOST environment variable")
	}
	if len(psTopic) <= 0 {
		log.Fatal("Failed to load PUBSUB_TOPIC environment variable")
	}
	if len(idLenght) <= 0 {
		log.Fatal("Failed to load ID_LENGHT environment variable")
	}
	if len(redisPass) <= 0 {
		log.Info("Using an empty Redis password")
	}

	parsedIdLenght, err := strconv.Atoi(idLenght)
	if err != nil {
		log.Fatal("Failed to parse ID_LENGHT environment variable")
	}

	defaulTTL := 60
	ttl, err := strconv.Atoi(redisTTL)
	if err != nil {
		ttl = defaulTTL
		log.Info("Using default Redis TTL: %v. Cause: %s", defaulTTL, err)
	} else {
		log.Info("Using Redis TTL: %v", ttl)
	}

	return EnvConfig{
		ProjectId:   project,
		RedisHost:   redisHost,
		RedisPass:   redisPass,
		RedisTTL:    ttl,
		PubsubTopic: psTopic,
		IdLength:    parsedIdLenght,
	}
}
