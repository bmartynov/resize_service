package util

import (
	"go.uber.org/zap"

	"github.com/bmartynov/resizer_service/crawler"
	resizeHandler "github.com/bmartynov/resizer_service/handler/resizer"
	"github.com/bmartynov/resizer_service/resizer"
	resizerSvc "github.com/bmartynov/resizer_service/service/resizer"
)

func NewLogger() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	return logger.Sugar()
}

func NewCrawler(config *crawler.Config) crawler.Crawler {
	return crawler.NewCrawler(*config)
}

func NewResizer(config *resizer.Config) resizer.Resizer {
	return resizer.NewJPEGResizer(*config)
}

func NewResizeHandler(logger *zap.SugaredLogger, svc resizerSvc.Resizer) resizeHandler.ResizeHandler {
	return resizeHandler.New(logger, svc)
}

func NewResizerService(
	logger *zap.SugaredLogger,
	crawler crawler.Crawler,
	resizer resizer.Resizer,
) resizerSvc.Resizer {
	return resizerSvc.New(logger, crawler, resizer)
}
