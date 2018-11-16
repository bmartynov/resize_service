package crawler

import (
	"context"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Crawler interface {
	Get(ctx context.Context, url string) (io.ReadCloser, error)
}

type httpCrawler struct {
	httpClient *http.Client
}

func (c *httpCrawler) Get(ctx context.Context, url string) (io.ReadCloser, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http.new_request")
	}

	request = request.WithContext(ctx)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http_client.do")
	}

	return response.Body, nil
}

func NewCrawler(config Config) Crawler {
	return &httpCrawler{
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}
