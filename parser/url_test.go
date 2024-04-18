package parser

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

var getToken = rapid.Custom(func(t *rapid.T) string {
	prefixes := []string{"http://", "https://", "ftp://", "file://", ""}

	rand.New(rand.NewSource(time.Now().Unix())) // initialize global pseudo random generator
	prefix := prefixes[rand.Intn(len(prefixes))]

	f := "<a href=\"%s%s.com\">"
	if time.Now().Unix()%2 == 0 {
		return fmt.Sprintf(f, prefix, rapid.StringMatching(
			"(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[1-9])\\.)"+
				"(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\\.)"+
				"{2}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])",
		).Draw(t, "ipv4"))
	}
	return fmt.Sprintf(f, prefix, rapid.StringMatching(
		"(^[a-zA-Z0-9]*){2,6}").Draw(t, "str"),
	)
})

func TestExtractURL(t *testing.T) {
	tests := map[string]func(t2 *testing.T){
		"check urls": rapid.MakeCheck(func(t *rapid.T) {
			token := getToken.Draw(t, "html")

			url, err := extractURL(token)
			if err != nil {
				assert.ErrorIs(t, err, errNoURLFound)
			}

			assert.Contains(t, token, url)
		}),
		"empty url": rapid.MakeCheck(func(t *rapid.T) {
			_, err := extractURL("blob")
			assert.ErrorIs(t, err, errNoURLFound)
		}),
	}

	for name, f := range tests {
		test := f
		t.Run(name, test)
	}
}
