package crawler

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpCrawler_Get(t *testing.T) {
	const payload = "1234567890"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, payload)
	}))
	defer s.Close()

	c := NewCrawler(Config{
		Timeout: time.Second,
	})

	r, err := c.Get(context.Background(), s.URL)
	if err != nil {
		t.Fatal(err)
	}

	rawPaylad, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	if string(rawPaylad) != payload {
		t.Fatal("payload mismatch")
	}
}

func TestHttpCrawler_Get_TimedOut(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
	}))
	defer s.Close()

	c := NewCrawler(Config{
		Timeout: time.Second,
	})

	_, err := c.Get(context.Background(), s.URL)
	if err == nil {
		t.Fatal("must be timeout")
	}
}
