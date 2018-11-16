package resizer

import (
	"net/url"
	"reflect"
	"testing"
)

func TestResizeRequestFrom(t *testing.T) {
	for _, tc := range []struct {
		name        string
		in          url.Values
		expected    *resizeRequest
		expectError error
	}{
		{
			name:     "ok",
			in:       url.Values{requestKeyUrl: []string{"http://ya.ru"}, requestKeySize: []string{"100x100"}},
			expected: &resizeRequest{Url: "http://ya.ru", Height: 100, Width: 100},
		},
		{
			name:        "invalid size",
			in:          url.Values{requestKeyUrl: []string{"http://ya.ru"}, requestKeySize: []string{"lolkek"}},
			expectError: errSizeInvalid,
		},
		{
			name:        "invalid url",
			in:          url.Values{requestKeyUrl: []string{"ya.ru"}, requestKeySize: []string{"100x100"}},
			expectError: errUrlInvalid,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, err := requestFrom(tc.in)
			if err != tc.expectError {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(r, tc.expected) {
				t.Fatalf("got != expected: got: %#v, expected: %#v", r, tc.expected)
			}
		})
	}
}
