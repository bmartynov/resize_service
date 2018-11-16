package main

import (
	"context"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/bmartynov/resizer_service/config"
	"github.com/bmartynov/resizer_service/handler/resizer"
	"github.com/bmartynov/resizer_service/util"
)

func run(
	lc fx.Lifecycle,
	logger *zap.SugaredLogger,
	handler resizer.ResizeHandler,
	config *config.HttpTransportConfig,
) {
	s := http.Server{
		Addr:    config.Address,
		Handler: handler.ResizeHandler(),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("start...")
			defer logger.Info("started.")

			go s.ListenAndServe()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("stop...")
			defer logger.Info("stopped.")

			return s.Shutdown(ctx)
		},
	})
}

func main() {
	logger := util.NewLogger()

	app := fx.New(
		fx.Provide(func() *zap.SugaredLogger {
			return logger
		}),
		fx.Provide(util.NewResizerConfig),
		fx.Provide(util.NewCrawler),
		fx.Provide(util.NewResizer),
		fx.Provide(util.NewResizerService),
		fx.Provide(util.NewResizeHandler),
		fx.Invoke(run),
	)
	app.Run()
}
