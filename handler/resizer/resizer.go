package resizer

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/bmartynov/resizer_service/service/resizer"
)

type ResizeHandler interface {
	ResizeHandler() http.Handler
}

type resizeHandler struct {
	logger *zap.SugaredLogger
	svc    resizer.Resizer
}

func (s *resizeHandler) handleResize(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.With("method", "handle_resize")

	logger.Info("start")

	query := r.URL.Query()

	logger = logger.With("query", query)

	rRequest, err := requestFrom(query)
	if err != nil {
		logger.Errorw("requestFrom", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tStart := time.Now()

	defer logger.Infow("stop", "took", time.Now().Sub(tStart))

	err = s.svc.Resize(r.Context(),
		rRequest.Url,
		rRequest.Width,
		rRequest.Height, w)

	if err != nil {
		logger.Errorw("svc.resize", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

func (s *resizeHandler) ResizeHandler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/resize/", s.handleResize).Methods(http.MethodGet)

	return r
}

func New(logger *zap.SugaredLogger, service resizer.Resizer) ResizeHandler {
	return &resizeHandler{
		logger: logger,
		svc:    service,
	}
}
