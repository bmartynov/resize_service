package resizer

import (
	"errors"
	"fmt"
	"github.com/bmartynov/resizer_service/cache"
	"go.uber.org/zap"
	"hash/fnv"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	requestKeyUrl  = "url"
	requestKeySize = "size"
)

var (
	errUrlInvalid  = errors.New("url invalid")
	errSizeInvalid = errors.New("size invalid")
)

type resizeRequest struct {
	Url    string
	Width  int
	Height int
}

func (r *resizeRequest) cacheKey() string {
	h := fnv.New64()
	h.Write([]byte(fmt.Sprintf("%s%d%d", r.Url, r.Width, r.Height)))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (r *resizeRequest) Validate() error {
	if r.Height == 0 || r.Width == 0 {
		return errSizeInvalid
	}

	u, err := url.Parse(r.Url)
	if err != nil {
		return errUrlInvalid
	}

	if u.Scheme == "" {
		return errUrlInvalid
	}

	return nil
}

func requestFrom(raw url.Values) (rRequest *resizeRequest, err error) {
	rawUrl := raw.Get(requestKeyUrl)
	if rawUrl == "" {
		return nil, errUrlInvalid
	}

	sizeParts := strings.SplitN(raw.Get(requestKeySize), "x", 2)
	if len(sizeParts) != 2 {
		return nil, errSizeInvalid
	}

	height, err := strconv.Atoi(sizeParts[0])
	if err != nil {
		return nil, errSizeInvalid
	}

	width, err := strconv.Atoi(sizeParts[1])
	if err != nil {
		return nil, errSizeInvalid
	}

	rRequest = &resizeRequest{
		Url:    rawUrl,
		Height: height,
		Width:  width,
	}

	if err = rRequest.Validate(); err != nil {
		return nil, err
	}

	return
}

type handlerFn func(w http.ResponseWriter, r *http.Request, rr *resizeRequest)

func withCache(
	cacheManager cache.Manager,
	logger *zap.SugaredLogger,
	fn handlerFn,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		rr, err := requestFrom(query)
		if err != nil {
			logger.Errorw("requestFrom", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cached, err := cacheManager.Get(rr.cacheKey())
		if err != nil {

			fn(w, r, rr)
		} else {
			_, err = io.Copy(w, cached)
			if err != nil {

			}
		}
	}
}
