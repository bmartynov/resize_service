package resizer

import (
	"context"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/bmartynov/resizer_service/cache"
	"github.com/bmartynov/resizer_service/service/resizer"
)

type ResizeHandler interface {
	ResizeHandler() http.Handler
}

type resizeHandler struct {
	logger *zap.SugaredLogger
	svc    resizer.Resizer
	cache  cache.Manager
}

func (s *resizeHandler) handleResize(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.With("method", "handle_resize")

	logger.Info("start")

	query := r.URL.Query()

	logger = logger.With("query", query)

	rr, err := requestFrom(query)
	if err != nil {
		logger.Errorw("requestFrom", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cacheKey := rr.cacheKey()

	cached, err := s.cache.Get(cacheKey)
	if err != nil {
		if cache.IsNoCacheError(err) {
			logger.Debugw("cache.get", "error", err, "cache", "miss")
		} else {
			logger.Errorw("cache.get", "error", err)
		}

		// this block used for parallel writing to cache and response without anesessary overhead
		wg := sync.WaitGroup{}

		pr, pw := io.Pipe()

		// tee is a Reader that writes to w what it reads from pr
		// tee used for store in cache, pr is resized image, w is response
		tee := io.TeeReader(pr, w)

		go func(wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()

			err = s.svc.Resize(context.Background(), rr.Url, rr.Width, rr.Height, pw)

			pw.CloseWithError(err)
		}(&wg)

		err = s.cache.Add(cacheKey, tee)
		if err != nil {
			logger.Errorw("cache.add", "error", err)
		}

		wg.Wait()
	} else {
		_, err := io.Copy(w, cached)
		// todo: handle headers
		if err != nil {
			logger.Errorw("response.write", "error", err)
		}
	}
}

func (s *resizeHandler) ResizeHandler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/resize/", s.handleResize).Methods(http.MethodGet)

	return r
}

func New(logger *zap.SugaredLogger, service resizer.Resizer, cache cache.Manager) ResizeHandler {
	return &resizeHandler{
		cache:  cache,
		logger: logger,
		svc:    service,
	}
}
