package helpers

import (
	"fmt"
	"net/url"
)

func IsUrl(str string) bool {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	return u.Scheme != "" && u.Host != ""
}

func GetFileUrl(fileId string, token string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s",
		token, url.QueryEscape(fileId))
}
