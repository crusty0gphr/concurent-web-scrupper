package parser

import (
	"io"

	"golang.org/x/net/html"
)

func ExtractValueByAttrName(body io.Reader, tag, attr string) (output []string) {
	tmp := make(map[string]struct{})
	t := html.NewTokenizer(body)

tokenParser:
	for {
		tokenType := t.Next()
		switch tokenType {
		case html.StartTagToken, html.EndTagToken:
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
		default:
			break tokenParser
		}
	}
	for key := range tmp {
		output = append(output, key)
	}
	return
}
