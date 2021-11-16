package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"ehgm.com.br/url-shortener/domain/model"
	"ehgm.com.br/url-shortener/domain/ports"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware(log ports.Logger) gin.HandlerFunc {
	return func(gc *gin.Context) {
		// before request
		gc.Next()
		// after request

		err := gc.Errors.Last()
		if err != nil {
			obJson := ErrorResponse{Message: err.Error(), Timestamp: time.Now()}
			log.Error("Request error: %s. Cause: %s", gc.Request.URL, err)

			var notFound *model.DocumentNotFoundError
			var invalidUrl *model.InvalidUrlError

			switch {
			case errors.As(err, &notFound):
				gc.JSON(http.StatusNotFound, obJson)
			case errors.As(err, &invalidUrl):
				gc.JSON(http.StatusBadRequest, obJson)
			default:
				gc.JSON(http.StatusInternalServerError, obJson)
			}
		}
	}
}

func validateUrl(rawUrl string) error {
	if len(rawUrl) > 2048 {
		return &model.InvalidUrlError{Messsage: "URL cannot be longer than 2048 characters"}
	}
	url, err := url.Parse(rawUrl)
	if err != nil || url.Scheme == "" || url.Host == "" {
		return &model.InvalidUrlError{Messsage: "URL does not have a valid format."}
	}
	return nil
}

func buildShortUrl(host, id string, isTLS bool) string {
	var url = fmt.Sprintf("https://%v/r/%v", host, id)
	if !isTLS {
		url = fmt.Sprintf("http://%v/r/%v", host, id)
	}
	return url
}
