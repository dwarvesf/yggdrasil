package toolkit

import (
	"net/url"

	uuid "github.com/satori/go.uuid"
)

// GetQueryParamFromURL return a value of given fieldName in param url
func GetQueryParamFromURL(URL url.URL, fieldName string) string {
	var queryValues = URL.Query()
	paramStr := queryValues.Get(fieldName)
	return paramStr
}

//IsUUIDZero check if uuid is zero value or not
func IsUUIDZero(u *uuid.UUID) bool {
	if u == nil {
		return true
	}
	for _, c := range u {
		if c != 0 {
			return false
		}
	}
	return true
}
