package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ehgm.com.br/url-shortener/domain/model"
	"ehgm.com.br/url-shortener/domain/ports"

	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/firestore"

	"github.com/go-redis/redis/v8"
)

var urlCollection = "urls"

// Struct that implements 'UrlRepository' interface
type urlRepository struct {
	log      ports.Logger
	fdb      *firestore.Client
	rdb      *redis.Client
	cacheTTL int
}

// Get an instance of 'UrlRepository' using this method
func NewUrlRepository(log ports.Logger,
	fdb *firestore.Client,
	rdb *redis.Client,
	cacheTTL int) ports.UrlRepository {

	return &urlRepository{log: log, fdb: fdb, rdb: rdb, cacheTTL: cacheTTL}
}

func (r *urlRepository) Save(ctx context.Context, id, url string, enable bool) error {
	shortUrl := model.ShortUrl{
		Id:     id,
		Url:    url,
		Enable: enable,
		Clicks: 0,
	}

	docRef, err := r.fdb.Collection(urlCollection).Doc(id).Create(ctx, shortUrl)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return &model.DocumentAlreadyExistsError{Id: id, Url: url}
		}
		return fmt.Errorf("Firestore creation error. %w", err)
	}

	shortUrl.CreateTime = docRef.UpdateTime
	go r.putInCache(&shortUrl)
	return err
}

func (r *urlRepository) FindById(ctx context.Context, id string) (*model.ShortUrl, error) {
	var shortUrl *model.ShortUrl
	var err error

	// Trying find in cache
	shortUrl = r.getFromCache(ctx, id)
	if *shortUrl == (model.ShortUrl{}) {

		// Trying find in NoSQL
		shortUrl, err = r.getFromNoSQL(ctx, id)
		if err != nil {
			return shortUrl, fmt.Errorf("FindById error. %w", err)
		}
		go r.putInCache(shortUrl)

	} else {
		r.log.Info("Id is cached: %v", id)
	}
	return shortUrl, nil
}

func (r *urlRepository) Update(ctx context.Context, id string, json map[string]interface{}) error {
	fields := []firestore.Update{}

	// Get the 2 allowed fields that can be updated
	for k, v := range json {
		if strings.EqualFold(k, "url") {
			fields = append(fields, firestore.Update{Path: k, Value: v})
		}
		if strings.EqualFold(k, "enable") {
			fields = append(fields, firestore.Update{Path: k, Value: v})
		}
	}
	if len(fields) <= 0 {
		r.log.Info("No attribute to update to Id: %v", id)
		return nil
	}

	_, err := r.fdb.Collection(urlCollection).Doc(id).Update(ctx, fields)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return &model.DocumentNotFoundError{Id: id}
		}
		return fmt.Errorf("Update Id error. %w", err)
	}
	r.log.Info("Id updated: %v", id)

	// Get from NoSQL to put in cache all fields
	go r.updateCache(id)
	return nil
}

func (r *urlRepository) GetStats(ctx context.Context, limit int) ([]model.ShortUrl, error) {
	shortUrls := []model.ShortUrl{}

	iter := r.fdb.Collection(urlCollection).OrderBy("clicks", firestore.Desc).Limit(limit).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return shortUrls, fmt.Errorf("GetStats error on %v element. %w", len(shortUrls), err)
		}
		temp := model.ShortUrl{}
		doc.DataTo(&temp)
		temp.CreateTime = doc.CreateTime
		shortUrls = append(shortUrls, []model.ShortUrl{temp}...)
	}

	r.log.Info("GetStats found %v urls", len(shortUrls))
	return shortUrls, nil
}

func (r *urlRepository) getFromNoSQL(ctx context.Context, id string) (*model.ShortUrl, error) {
	var shortUrl model.ShortUrl

	dsnap, err := r.fdb.Collection(urlCollection).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return &shortUrl, &model.DocumentNotFoundError{Id: id}
		}
		return &shortUrl, fmt.Errorf("getFromNoSQL error. %w", err)
	}

	dsnap.DataTo(&shortUrl)
	shortUrl.CreateTime = dsnap.CreateTime
	return &shortUrl, nil
}

func (r *urlRepository) getFromCache(ctx context.Context, id string) *model.ShortUrl {
	var shortUrl model.ShortUrl

	text, err := r.rdb.Get(ctx, id).Result()
	if err != nil {
		if err == redis.Nil {
			r.log.Info("Id not found in cache: %v", id)
			return &shortUrl
		}
		r.log.Error("getFromCache error for Id: %v. Cause: %s", id, err)
		return &shortUrl
	}

	shortUrl, err = jsonToStruct(text)
	if err != nil {
		r.log.Error("getFromCache error on jsonToStruct for Id: %v. Cause: %s", id, err)
	}
	return &shortUrl
}

func (r *urlRepository) putInCache(shortUrl *model.ShortUrl) {
	id := shortUrl.Id
	duration := time.Duration(r.cacheTTL) * time.Minute

	value, err := structToJson(shortUrl)
	if err != nil {
		r.log.Error("structToJson error for Id: %v. Cause: %s", id, err)
		return
	}

	ctx := context.Background()
	err = r.rdb.Set(ctx, id, value, duration).Err()
	if err != nil {
		r.log.Error("putInCache error for Id: %v. Cause: %s", id, err)
	} else {
		r.log.Info("Put Id in cache %v: ", id)
	}
}

func (r *urlRepository) updateCache(id string) {
	ctx := context.Background()
	shortUrl, err := r.getFromNoSQL(ctx, id)
	if err != nil {
		r.log.Error("updateCache error after Update NoSQL. Id: %v. Cause: %s", id, err)
	}
	r.putInCache(shortUrl)
}
