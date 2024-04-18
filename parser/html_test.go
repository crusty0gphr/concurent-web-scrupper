package parser

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

const (
	anchor = "a"
	size   = math.MaxUint8
	ipv4   = "(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[1-9])\\.)" +
		"(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\\.)" +
		"{2}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])"
	url = "(^[a-zA-Z0-9]*){6}"
)

var genToken = rapid.Custom(func(t *rapid.T) []string {
	format := "<a href=\"%s%s.com\">"
	prefixes := []string{"http://", "https://", "ftp://", "file://", ""}

	rand.New(rand.NewSource(time.Now().Unix())) // initialize global pseudo random generator
	prefix := prefixes[rand.Intn(len(prefixes))]

	out := make([]string, size)
	for i := range out {
		if time.Now().Unix()%2 == 0 {
			out[i] = fmt.Sprintf(format, prefix, rapid.StringMatching(ipv4).Filter(func(s string) bool {
				return len(s) > 0
			}).Draw(t, "ip"))
		}
		out[i] = fmt.Sprintf(format, prefix, rapid.StringMatching(url).Filter(func(s string) bool {
			return len(s) > 0
		}).Draw(t, "text"))
	}

	return out
})

func TestExtractValueByAttrName(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		tokens := genToken.Draw(t, "html")
		_ = tokens

		var b strings.Builder
		b.WriteString("<body>")
		b.WriteString(strings.Join(tokens[:], ""))
		b.WriteString("</body>")

		res := ExtractValueByAttrName(strings.NewReader(b.String()), anchor, "href")
		assert.Len(t, res, size)
	})
}
