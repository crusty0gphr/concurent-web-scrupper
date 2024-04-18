package parser

import (
	"errors"
	"regexp"
)

var errNoURLFound = errors.New("no urls in string")

const regexMatch = `([\w+]+\:\/\/)?([\w\d-]+\.)*[\w-]+[\.\:]\w+([\/\?\=\&\#\.]?[\w-]+)*\/?`

func extractURL(t string) (string, error) {
	compile, err := regexp.Compile(regexMatch)
	if err != nil {
		return "", err
	}

	r := compile.FindString(t)
	if len(r) == 0 {
		return "", errNoURLFound
	}

	return r, nil
}
