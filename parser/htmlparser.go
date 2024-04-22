package parser

import (
	"io"

	"golang.org/x/net/html"
)

func ExtractValueByAttrName(body io.Reader, tag, attr string) []string {
	tmp := make(map[string]struct{})
	t := html.NewTokenizer(body)

tokenParser:
	for {
		tokenType := t.Next()
		if tokenType == html.ErrorToken {
			break tokenParser
		}

		token := t.Token()
		if token.Data != tag {
			continue
		}
		for _, attribute := range token.Attr {
			if attribute.Key != attr {
				continue
			}
			// avoid duplicates
			if _, ok := tmp[attribute.Val]; !ok {
				tmp[attribute.Val] = struct{}{} // saving little space here, for no reason...
			}
		}
	}
	var out []string
	for key := range tmp {
		out = append(out, key)
	}
	return out
}
