package resizer

import (
	"context"
	"io"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/bmartynov/resizer_service/crawler"
	"github.com/bmartynov/resizer_service/resizer"
)

type Resizer interface {
	Resize(ctx context.Context, url string, height, width int, dst io.Writer) (err error)
}

type Service struct {
	logger  *zap.SugaredLogger
	crawler crawler.Crawler
	resizer resizer.Resizer
}

func (s *Service) Resize(ctx context.Context, url string, height, width int, dst io.Writer) (err error) {
	tStart := time.Now()

	logger := s.logger.With("method", "resize", "url", url, "height", height, "width", width)

	logger.Info("start")
	defer s.logger.Infow("stop", "took", time.Now().Sub(tStart))

	src, err := s.crawler.Get(ctx, url)
	if err != nil {
		return errors.Wrap(err, "crawler.get")
	}
	defer src.Close()

	err = s.resizer.Resize(src, width, height, dst)
	if err != nil {
		return errors.Wrap(err, "resizer.resize")
	}

	return nil
}

func New(
	logger *zap.SugaredLogger,
	crawler crawler.Crawler,
	resizer resizer.Resizer,
) Resizer {
	return &Service{
		logger:  logger,
		crawler: crawler,
		resizer: resizer,
	}
}
