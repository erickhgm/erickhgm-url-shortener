package main

import (
	"context"

	"ehgm.com.br/url-shortener/adapters/api"
	"ehgm.com.br/url-shortener/adapters/idgenerator"
	"ehgm.com.br/url-shortener/adapters/pubsub"
	"ehgm.com.br/url-shortener/adapters/repository"
	"ehgm.com.br/url-shortener/config"
	"ehgm.com.br/url-shortener/domain/ports"
	"ehgm.com.br/url-shortener/domain/usecases"

	"github.com/gin-gonic/gin"
)

var log ports.Logger

func init() {
	log = config.NewLogger()
	log.Info("URL Shortener current version: 1.0.0")
}

func main() {
	log.Info("Loading dependencies ...")

	ctx := context.Background()
	env := config.NewEnvConfig(log)
	rdb := config.NewRedisClient(env.RedisHost, env.RedisPass)
	ps := config.NewPubSubClient(ctx, log, env.ProjectId)
	fdb := config.NewFirestoreClient(ctx, log, env.ProjectId)

	idGenerator := idgenerator.NewIdGenerator(env.IdLength)
	urlCounter := pubsub.NewUrlCounter(log, ps, env.PubsubTopic)
	urlRepository := repository.NewUrlRepository(log, fdb, rdb, env.RedisTTL)
	urlService := usecases.NewUrlService(log, idGenerator, urlRepository, urlCounter)
	controller := api.NewUrlController(log, urlService)

	log.Info("Starting Gin server ...")

	router := gin.Default()
	router.Use(api.ErrorHandlerMiddleware(log))

	docGroup := router.Group("/doc")
	docGroup.Static("/", "./doc")

	docStatic := router.Group("/static")
	docStatic.Static("/", "./static")

	redirectGroup := router.Group("/r")
	redirectGroup.GET("/:id", controller.RedirectToUrl)

	urlsGroup := router.Group("/urls")
	urlsGroup.POST("/", controller.PostUrl)
	urlsGroup.GET("/:id", controller.GetUrl)
	urlsGroup.PATCH("/:id", controller.PatchUrl)

	statsGroup := router.Group("/stats")
	statsGroup.GET("/", controller.GetStats)

	router.Run()
}
