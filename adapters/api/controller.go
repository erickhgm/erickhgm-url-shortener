package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"ehgm.com.br/url-shortener/domain/model"
	"ehgm.com.br/url-shortener/domain/ports"

	"github.com/gin-gonic/gin"
)

// This struct does not need an interface, its the first level of dependency injection
type urlController struct {
	log        ports.Logger
	urlService ports.UrlService
}

// Get an instance of 'urlController' using this method
func NewUrlController(log ports.Logger, urlService ports.UrlService) *urlController {
	return &urlController{log: log, urlService: urlService}
}

func (c *urlController) PostUrl(gc *gin.Context) {
	var json Url
	ctx := gc.Request.Context()

	if err := gc.BindJSON(&json); err != nil {
		gc.Error(fmt.Errorf("BindJSON error in urlService.PostUrl. %w", err))
		return
	}

	if err := validateUrl(json.Url); err != nil {
		gc.Error(fmt.Errorf("validateUrl error in urlService.PostUrl. %w", err))
		return
	}

	id, err := c.urlService.GenerateId(ctx, json.Url)
	if err != nil {
		gc.Error(fmt.Errorf("GenerateId error in urlService.PostUrl. %w", err))
		return
	}

	var url = buildShortUrl(gc.Request.Host, id, gc.Request.TLS != nil)
	c.log.Info("Long Url: %v and short Url: %v", json.Url, url)

	gc.JSON(http.StatusCreated, Url{Url: url})
}

func (c *urlController) RedirectToUrl(gc *gin.Context) {
	var err error
	ctx := gc.Request.Context()

	id := gc.Param("id")
	url, enable, err := c.urlService.GetUrlToRedirect(ctx, id)
	if err != nil {
		gc.Error(fmt.Errorf("GetUrlToRedirect error in urlService.RedirectToUrl. %w", err))
		return
	}

	switch {
	case len(url) <= 0:
		gc.Status(http.StatusNotFound)
	case !enable:
		gc.Redirect(http.StatusFound, "/static/423.html")
	default:
		gc.Redirect(http.StatusFound, url)
	}
}

func (c *urlController) PatchUrl(gc *gin.Context) {
	var err error
	ctx := gc.Request.Context()

	defer gc.Request.Body.Close()
	bodyBytes, err := ioutil.ReadAll(gc.Request.Body)
	if err != nil {
		gc.Error(fmt.Errorf("Read Body error in urlService.PatchUrl. %w", err))
		return
	}

	var jsonBody map[string]interface{}
	if err = json.Unmarshal([]byte(bodyBytes), &jsonBody); err != nil {
		gc.Error(fmt.Errorf("Unmarshal error in urlService.PatchUrl. %w", err))
		return
	}

	id := gc.Param("id")
	if err = c.urlService.UpdateUrl(ctx, id, jsonBody); err != nil {
		gc.Error(fmt.Errorf("UpdateUrl error in urlService.PatchUrl. %w", err))
		return
	}

	gc.Status(http.StatusOK)
}

func (c *urlController) GetUrl(gc *gin.Context) {
	ctx := gc.Request.Context()

	id := gc.Param("id")
	shortUrl, err := c.urlService.GetUrl(ctx, id)
	if err != nil {
		gc.Error(fmt.Errorf("GetUrl error in urlService.GetUrl. %w", err))
		return
	}

	if *shortUrl == (model.ShortUrl{}) {
		gc.Status(http.StatusNotFound)
	} else {
		gc.JSON(http.StatusOK, shortUrl)
	}
}

func (c *urlController) GetStats(gc *gin.Context) {
	ctx := gc.Request.Context()

	lim := 0
	if limit, ok := gc.GetQuery("limit"); ok {
		v, _ := strconv.Atoi(limit)
		lim = v
	}

	urls, err := c.urlService.GetStats(ctx, lim)
	if err != nil {
		gc.Error(fmt.Errorf("GetStats error in urlService.GetStats. %w", err))
		return
	}

	gc.JSON(http.StatusOK, urls)
}
