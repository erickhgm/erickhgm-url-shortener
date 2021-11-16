package usecases

import (
	"context"
	"errors"
	"testing"

	"ehgm.com.br/url-shortener/domain/model"
	"ehgm.com.br/url-shortener/domain/ports"
)

// Empty UrlRepository
type urlRepositoryMock struct {
	saveFn     func(ctx context.Context, id, url string, enable bool) error
	findByIdFn func(ctx context.Context, id string) (*model.ShortUrl, error)
	updateFn   func(ctx context.Context, id string, json map[string]interface{}) error
	getStatsFn func(ctx context.Context, limit int) ([]model.ShortUrl, error)
}

func (r *urlRepositoryMock) Save(ctx context.Context, id, url string, enable bool) error {
	if r.saveFn != nil {
		return r.saveFn(ctx, id, url, enable)
	}
	return nil
}

func (r *urlRepositoryMock) FindById(ctx context.Context, id string) (*model.ShortUrl, error) {
	if r.findByIdFn != nil {
		return r.findByIdFn(ctx, id)
	}
	return &model.ShortUrl{}, nil
}

func (r *urlRepositoryMock) Update(ctx context.Context, id string, json map[string]interface{}) error {
	if r.updateFn != nil {
		return r.updateFn(ctx, id, json)
	}
	return nil
}

func (r *urlRepositoryMock) GetStats(ctx context.Context, limit int) ([]model.ShortUrl, error) {
	if r.getStatsFn != nil {
		return r.getStatsFn(ctx, limit)
	}
	return []model.ShortUrl{}, nil
}

// Empty IdGenerator
type idGeneratorMock struct {
	newFn func() (string, error)
}

func (g *idGeneratorMock) New() (string, error) {
	if g.newFn != nil {
		return g.newFn()
	}
	return "", nil
}

// Empty UrlCounter
type urlCounterMock struct{}

func (c *urlCounterMock) IncrementCounter(id string) {}

// Empty Logger
type loggerMock struct{}

func (l *loggerMock) Info(format string, v ...interface{})  {}
func (l *loggerMock) Error(format string, v ...interface{}) {}
func (l *loggerMock) Fatal(format string, v ...interface{}) {}

func TestGenerateId(t *testing.T) {
	type Input struct {
		log         ports.Logger
		idGenerator ports.IdGenerator
		urlCounter  ports.UrlCounter
		repo        ports.UrlRepository
		url         string
		counter     int
	}

	type Output struct {
		id       string
		hasError bool
	}

	tests := map[string]struct {
		input  Input
		output Output
	}{
		"Test 01 - Should return a new valid Id": {
			Input{
				log: &loggerMock{},
				idGenerator: &idGeneratorMock{
					newFn: func() (string, error) {
						return "1q2w3e", nil
					},
				},
				urlCounter: &urlCounterMock{},
				repo: &urlRepositoryMock{
					saveFn: func(ctx context.Context, id, url string, enable bool) error {
						return nil
					}},
				url: "https://ehgm.com.br"},
			Output{
				id:       "1q2w3e",
				hasError: false,
			}},

		"Test 02 - Should return a NanoId error": {
			Input{
				log: &loggerMock{},
				idGenerator: &idGeneratorMock{
					newFn: func() (string, error) {
						return "", errors.New("NanoId error")
					},
				},
				urlCounter: &urlCounterMock{},
				repo:       &urlRepositoryMock{},
				url:        "https://ehgm.com.br"},
			Output{
				id:       "",
				hasError: true,
			}},

		"Test 03 - Should return a Save error": {
			Input{
				log:         &loggerMock{},
				idGenerator: &idGeneratorMock{},
				urlCounter:  &urlCounterMock{},
				repo: &urlRepositoryMock{
					saveFn: func(ctx context.Context, id, url string, enable bool) error {
						return errors.New("Save error")
					}},
				url: "https://ehgm.com.br"},
			Output{
				id:       "",
				hasError: true,
			}},
	}

	ctx := context.Background()

	for i, test := range tests {
		urlService := NewUrlService(test.input.log, test.input.idGenerator, test.input.repo, test.input.urlCounter)
		id, err := urlService.GenerateId(ctx, test.input.url)

		if test.output.hasError && err == nil {
			t.Errorf("#%s: Output is: %s. But should has error: %v", i, err, test.output.hasError)
			continue
		}
		if !test.output.hasError && err != nil {
			t.Errorf("#%s: Output is: %s. But should not has error: %v", i, err, test.output.hasError)
			continue
		}
		if !test.output.hasError && id != test.output.id {
			t.Errorf("#%s: Output is: %v. But should be: %v", i, id, test.output.id)
		}
	}
}

func TestGenerateIdDocumentAlrearyExist(t *testing.T) {
	type Input struct {
		log         ports.Logger
		idGenerator ports.IdGenerator
		urlCounter  ports.UrlCounter
		repo        ports.UrlRepository
		url         string
		counter     int
	}

	type Output struct {
		id       string
		hasError bool
	}

	var saveCounter int

	tests := map[string]struct {
		input  Input
		output Output
	}{
		"Test 01 - Should return a new valid id": {
			Input{
				log: &loggerMock{},
				idGenerator: &idGeneratorMock{
					newFn: func() (string, error) {
						return "1q2w3e", nil
					},
				},
				urlCounter: &urlCounterMock{},
				repo: &urlRepositoryMock{
					saveFn: func(ctx context.Context, id, url string, enable bool) error {
						saveCounter++
						if saveCounter > 1 {
							return nil
						}
						return &model.DocumentAlreadyExistsError{}
					}},
				url: "https://ehgm.com.br"},
			Output{
				id:       "1q2w3e",
				hasError: false,
			}},

		"Test 02 - Should return a DocumentAlreadyExistsError error": {
			Input{
				log:         &loggerMock{},
				idGenerator: &idGeneratorMock{},
				urlCounter:  &urlCounterMock{},
				repo: &urlRepositoryMock{
					saveFn: func(ctx context.Context, id, url string, enable bool) error {
						return &model.DocumentAlreadyExistsError{}
					}},
				url: "https://ehgm.com.br"},
			Output{
				id:       "1q2w3e",
				hasError: true,
			}},
	}

	ctx := context.Background()

	for i, test := range tests {
		urlService := NewUrlService(test.input.log, test.input.idGenerator, test.input.repo, test.input.urlCounter)
		id, err := urlService.GenerateId(ctx, test.input.url)

		if test.output.hasError && err == nil {
			t.Errorf("#%s: Output is: %s. But should has error: %v", i, err, test.output.hasError)
			continue
		}
		if !test.output.hasError && err != nil {
			t.Errorf("#%s: Output is: %s. But should not has error: %v", i, err, test.output.hasError)
			continue
		}
		if !test.output.hasError && id != test.output.id {
			t.Errorf("#%s: Output is: %v. But should be: %v", i, id, test.output.id)
		}
	}
}

func TestGetUrl(t *testing.T) {
	type Input struct {
		log         ports.Logger
		idGenerator ports.IdGenerator
		urlCounter  ports.UrlCounter
		repo        ports.UrlRepository
		id          string
	}

	type Output struct {
		shortUrl model.ShortUrl
		hasError bool
	}

	tests := map[string]struct {
		input  Input
		output Output
	}{
		"Test 01 - Should return a URL": {
			Input{
				log:         &loggerMock{},
				idGenerator: &idGeneratorMock{},
				urlCounter:  &urlCounterMock{},
				repo: &urlRepositoryMock{
					findByIdFn: func(ctx context.Context, id string) (*model.ShortUrl, error) {
						return &model.ShortUrl{Url: "https://ehgm.com.br", Enable: true}, nil
					}},
				id: "1q2w3e"},
			Output{
				shortUrl: model.ShortUrl{Url: "https://ehgm.com.br", Enable: true},
				hasError: false,
			}},

		"Test 02 - Should return error": {
			Input{
				log:         &loggerMock{},
				idGenerator: &idGeneratorMock{},
				urlCounter:  &urlCounterMock{},
				repo: &urlRepositoryMock{
					findByIdFn: func(ctx context.Context, id string) (*model.ShortUrl, error) {
						return &model.ShortUrl{}, errors.New("FindById error")
					}},
				id: "1q2w3e"},
			Output{
				shortUrl: model.ShortUrl{},
				hasError: true,
			}},
	}

	ctx := context.Background()

	for i, test := range tests {
		urlService := NewUrlService(test.input.log, test.input.idGenerator, test.input.repo, test.input.urlCounter)
		shortUrl, err := urlService.GetUrl(ctx, test.input.id)

		if test.output.hasError && err == nil {
			t.Errorf("#%s: Output is: %s. But should has error: %v", i, err, test.output.hasError)
			continue
		}
		if !test.output.hasError && err != nil {
			t.Errorf("#%s: Output is: %s. But should not has error: %v", i, err, test.output.hasError)
			continue
		}
		if !test.output.hasError && *shortUrl != test.output.shortUrl {
			t.Errorf("#%s: Output is: %v. But should be: %v", i, shortUrl, test.output.shortUrl)
		}
	}
}

func TestGetUrlToRedirect(t *testing.T) {
	type Input struct {
		log         ports.Logger
		idGenerator ports.IdGenerator
		urlCounter  ports.UrlCounter
		repo        ports.UrlRepository
		id          string
	}

	type Output struct {
		url      string
		enable   bool
		hasError bool
	}

	tests := map[string]struct {
		input  Input
		output Output
	}{
		"Test 01 - Should return an enabled URL": {
			Input{
				log:         &loggerMock{},
				idGenerator: &idGeneratorMock{},
				urlCounter:  &urlCounterMock{},
				repo: &urlRepositoryMock{
					findByIdFn: func(ctx context.Context, id string) (*model.ShortUrl, error) {
						return &model.ShortUrl{Url: "https://ehgm.com.br", Enable: true}, nil
					}},
				id: "1q2w3e"},
			Output{
				url:      "https://ehgm.com.br",
				enable:   true,
				hasError: false,
			}},

		"Test 02 - Should return an disabled URL": {
			Input{
				log:         &loggerMock{},
				idGenerator: &idGeneratorMock{},
				urlCounter:  &urlCounterMock{},
				repo: &urlRepositoryMock{
					findByIdFn: func(ctx context.Context, id string) (*model.ShortUrl, error) {
						return &model.ShortUrl{Url: "https://ehgm.com.br", Enable: false}, nil
					}},
				id: "1q2w3e"},
			Output{
				url:      "https://ehgm.com.br",
				enable:   false,
				hasError: false,
			}},

		"Test 03 - Should return error": {
			Input{
				log:         &loggerMock{},
				idGenerator: &idGeneratorMock{},
				urlCounter:  &urlCounterMock{},
				repo: &urlRepositoryMock{
					findByIdFn: func(ctx context.Context, id string) (*model.ShortUrl, error) {
						return &model.ShortUrl{}, errors.New("FindById error")
					}},
				id: "1q2w3e"},
			Output{
				url:      "https://ehgm.com.br",
				enable:   true,
				hasError: true,
			}},
	}

	ctx := context.Background()

	for i, test := range tests {
		urlService := NewUrlService(test.input.log, test.input.idGenerator, test.input.repo, test.input.urlCounter)
		url, enable, err := urlService.GetUrlToRedirect(ctx, test.input.id)

		if test.output.hasError && err == nil {
			t.Errorf("#%s: Output is: %s. But should has error: %v", i, err, test.output.hasError)
			continue
		}
		if !test.output.hasError && err != nil {
			t.Errorf("#%s: Output is: %s. But should not has error: %v", i, err, test.output.hasError)
			continue
		}
		if !test.output.hasError && (url != test.output.url || enable != test.output.enable) {
			t.Errorf("#%s: Output is: %v / %v. But should be: %v / %v", i, url, enable, test.output.url, test.output.enable)
		}
	}
}

func TestUpdateUrl(t *testing.T) {
	type Input struct {
		log         ports.Logger
		idGenerator ports.IdGenerator
		urlCounter  ports.UrlCounter
		repo        ports.UrlRepository
		id          string
		json        map[string]interface{}
	}

	type Output struct {
		hasError bool
	}

	tests := map[string]struct {
		input  Input
		output Output
	}{
		"Test 01 - Should call and return a nil error": {
			Input{
				log:         &loggerMock{},
				idGenerator: &idGeneratorMock{},
				urlCounter:  &urlCounterMock{},
				repo: &urlRepositoryMock{
					updateFn: func(ctx context.Context, id string, json map[string]interface{}) error {
						return nil
					}},
				id:   "1q2w3e",
				json: map[string]interface{}{},
			},
			Output{hasError: false},
		},

		"Test 02 - Should call and return an error": {
			Input{
				log:         &loggerMock{},
				idGenerator: &idGeneratorMock{},
				urlCounter:  &urlCounterMock{},
				repo: &urlRepositoryMock{
					updateFn: func(ctx context.Context, id string, json map[string]interface{}) error {
						return errors.New("Update error")
					}},
				id:   "1q2w3e",
				json: map[string]interface{}{},
			},
			Output{hasError: true},
		},
	}

	ctx := context.Background()

	for i, test := range tests {
		urlService := NewUrlService(test.input.log, test.input.idGenerator, test.input.repo, test.input.urlCounter)
		err := urlService.UpdateUrl(ctx, test.input.id, test.input.json)

		if test.output.hasError && err == nil {
			t.Errorf("#%s: Output is: %s. But should has error: %v", i, err, test.output.hasError)
		}
		if !test.output.hasError && err != nil {
			t.Errorf("#%s: Output is: %s. But should not has error: %v", i, err, test.output.hasError)
		}
	}
}

func TestGetStats(t *testing.T) {
	type Input struct {
		log         ports.Logger
		idGenerator ports.IdGenerator
		urlCounter  ports.UrlCounter
		repo        ports.UrlRepository
		limit       int
	}

	type Output struct {
		shortUrls []model.ShortUrl
		hasError  bool
	}

	tests := map[string]struct {
		input  Input
		output Output
	}{
		"Test 01 - Should call with limit=1": {
			Input{
				log:         &loggerMock{},
				idGenerator: &idGeneratorMock{},
				urlCounter:  &urlCounterMock{},
				repo: &urlRepositoryMock{
					getStatsFn: func(ctx context.Context, limit int) ([]model.ShortUrl, error) {
						if limit != 1 {
							return []model.ShortUrl{}, errors.New("Limit error")
						}
						return make([]model.ShortUrl, 1), nil
					}},
				limit: 1},
			Output{
				shortUrls: make([]model.ShortUrl, 1),
				hasError:  false,
			},
		},

		"Test 02 - Should call with limit=10": {
			Input{
				log:         &loggerMock{},
				idGenerator: &idGeneratorMock{},
				urlCounter:  &urlCounterMock{},
				repo: &urlRepositoryMock{
					getStatsFn: func(ctx context.Context, limit int) ([]model.ShortUrl, error) {
						if limit != 10 {
							return []model.ShortUrl{}, errors.New("Limit error")
						}
						return make([]model.ShortUrl, 10), nil
					}},
				limit: 0},
			Output{
				shortUrls: make([]model.ShortUrl, 10),
				hasError:  false,
			},
		},

		"Test 03 - Should call and return error": {
			Input{
				log:         &loggerMock{},
				idGenerator: &idGeneratorMock{},
				urlCounter:  &urlCounterMock{},
				repo: &urlRepositoryMock{
					getStatsFn: func(ctx context.Context, limit int) ([]model.ShortUrl, error) {
						return []model.ShortUrl{}, errors.New("GetStats error")
					}},
				limit: 100},
			Output{
				shortUrls: []model.ShortUrl{},
				hasError:  true,
			},
		},
	}

	ctx := context.Background()

	for i, test := range tests {
		urlService := NewUrlService(test.input.log, test.input.idGenerator, test.input.repo, test.input.urlCounter)
		shortUrls, err := urlService.GetStats(ctx, test.input.limit)

		if test.output.hasError && err == nil {
			t.Errorf("#%s: Output is: %s. But should has error: %v", i, err, test.output.hasError)
		}
		if !test.output.hasError && err != nil {
			t.Errorf("#%s: Output is: %s. But should not has error: %v", i, err, test.output.hasError)
		}
		if !test.output.hasError && len(shortUrls) != len(test.output.shortUrls) {
			t.Errorf("#%s: Output is: %v. But should be: %v", i, len(shortUrls), len(test.output.shortUrls))
		}
	}
}
