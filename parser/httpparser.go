package parser

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	anchor = "a"
	target = "href"
)

func ExtractUrls(url string) ([]string, error) {
	status, body := ping(url)
	if status >= http.StatusBadRequest {
		return nil, errors.New(fmt.Sprintf("status code %d", status))
	}

	urls := ExtractValueByAttrName(body, anchor, target)
	var out []string
	for _, u := range urls {
		if validateURL(u) {
			continue
		}
		out = append(out, u)
	}
	return out, nil
}

func validateURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func ping(url string) (int, io.ReadCloser) {
	var (
		status int
		body   io.ReadCloser
	)

	resp, err := http.Get(url)
	if err != nil {
		status = http.StatusBadRequest
	} else {
		status = resp.StatusCode
		body = resp.Body
	}
	return status, body
}
