package helpers

import (
	"net/url"
)

func IsUrl(str string) bool {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	return u.Scheme != "" && u.Host != ""
}
