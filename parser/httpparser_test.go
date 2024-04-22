package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateURl(t *testing.T) {
	tests := map[string]struct {
		url   string
		valid bool
	}{
		"valid": {
			url:   "https://asdf.com",
			valid: true,
		},
		"invalid": {
			url:   "asdf.html",
			valid: false,
		},
	}

	for name, s := range tests {
		test := s
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.valid, validateURL(test.url))
		})
	}
}

func TestHTTPParser(t *testing.T) {
	tests := map[string]func(t *testing.T){
		"example.com": func(t *testing.T) {
			res, err := ExtractUrls("https://example.com")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, "https://www.iana.org/domains/example", res[0])
		},
		"asdf.com": func(t *testing.T) {
			res, err := ExtractUrls("https://asdf.com")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, "//asdfforums.com", res[0])
		},
	}

	for name, s := range tests {
		test := s
		t.Run(name, test)
	}
}

func TestPing(t *testing.T) {
	tests := map[string]struct {
		url      string
		expected int
	}{
		"200": {
			url:      "https://httpstat.us/200",
			expected: 200,
		},
		"404": {
			url:      "https://httpstat.us/404",
			expected: 404,
		},
		"500": {
			url:      "https://httpstat.us/500",
			expected: 500,
		},
	}

	for name, s := range tests {
		test := s
		t.Run(name, func(t *testing.T) {
			status, _ := ping(test.url)
			assert.Equal(t, test.expected, status)
		})
	}
}
