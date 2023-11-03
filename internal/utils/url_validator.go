/*
Copyright © 2023 grarich <grarich@grawlily.com>
*/
package utils

import (
	"net/url"
)

func IsValidUrl(str string) bool {
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	return true
}
