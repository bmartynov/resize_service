package resizer

import (
	"bytes"
	"os"
	"testing"

	"github.com/pkg/errors"
)

func TestNewJPEGResizer(t *testing.T) {
	r := NewJPEGResizer(Config{Filter: "box"})

	original, err := os.Open("./assets/original.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer original.Close()

	b := bytes.NewBuffer(nil)

	err = r.Resize(original, 400, 400, b)
	if err != nil {
		t.Fatal(errors.Wrap(err, "must be resized"))
	}
}
