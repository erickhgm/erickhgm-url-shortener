package usecases

import (
	"context"
	"errors"
	"fmt"

	"ehgm.com.br/url-shortener/domain/model"
	"ehgm.com.br/url-shortener/domain/ports"
)

// Struct that implements 'UrlService' interface
type urlService struct {
	log           ports.Logger
	idGenerator   ports.IdGenerator
	urlRepository ports.UrlRepository
	urlCounter    ports.UrlCounter
}

// Get an instance of 'UrlService' using this method
func NewUrlService(log ports.Logger,
	idGenerator ports.IdGenerator,
	urlRepository ports.UrlRepository,
	urlCounter ports.UrlCounter) ports.UrlService {

	return &urlService{log: log, idGenerator: idGenerator, urlRepository: urlRepository, urlCounter: urlCounter}
}

func (s *urlService) GenerateId(ctx context.Context, url string) (string, error) {
	var id string
	var err error

	// If already exist, generate other id end try again
	// This will rarely happen, we have 4.398.046.511.104 different ids (4.3 Trillion)

	for i := 1; i <= 3; i++ {
		if id, err = s.idGenerator.New(); err != nil {
			return "", fmt.Errorf("Nano Id generation error. %w", err)
		}

		err = s.urlRepository.Save(ctx, id, url, true)
		if err != nil {
			var docExist *model.DocumentAlreadyExistsError

			// Already in use at database
			if errors.As(err, &docExist) {
				s.log.Info("Retrying generate Id for Url: %v. Num: %v", url, i)
				continue
			}
			return "", fmt.Errorf("Save Id %v error. %w", id, err)
		}

		s.log.Info("Successfully generated id: %v for Url: %v", id, url)
		break
	}
	return id, err
}

func (s *urlService) GetUrl(ctx context.Context, id string) (*model.ShortUrl, error) {
	shortUrl, err := s.urlRepository.FindById(ctx, id)
	if err != nil {
		return shortUrl, fmt.Errorf("GetUrl error for Id: %v. %w", id, err)
	}
	return shortUrl, nil
}

func (s *urlService) GetUrlToRedirect(ctx context.Context, id string) (string, bool, error) {
	shortUrl, err := s.urlRepository.FindById(ctx, id)
	if err != nil {
		return "", false, fmt.Errorf("GetUrlToRedirect error for Id: %v. %w", id, err)
	}

	var url string
	if *shortUrl != (model.ShortUrl{}) {
		url = shortUrl.Url
		go s.urlCounter.IncrementCounter(id)
	}
	return url, shortUrl.Enable, nil
}

func (s *urlService) UpdateUrl(ctx context.Context, id string, json map[string]interface{}) error {
	err := s.urlRepository.Update(ctx, id, json)
	if err != nil {
		return fmt.Errorf("UpdateUrl error for Id: %v. %w", id, err)
	}
	return nil
}

func (s *urlService) GetStats(ctx context.Context, limit int) ([]model.ShortUrl, error) {
	var defaultLimit = 10

	if limit > 0 {
		defaultLimit = limit
	} else {
		s.log.Info("Using limit default: %v, limit received: %v", defaultLimit, limit)
	}

	shortUrls, err := s.urlRepository.GetStats(ctx, defaultLimit)
	if err != nil {
		return nil, fmt.Errorf("GetStats error using limit: %v. %w", defaultLimit, err)
	}
	return shortUrls, nil
}
