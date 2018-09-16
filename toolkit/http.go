package toolkit

import "net/url"

// GetQueryParamFromURL return a value of given fieldName in param url
func GetQueryParamFromURL(URL url.URL, fieldName string) string {
	var queryValues = URL.Query()
	paramStr := queryValues.Get(fieldName)
	return paramStr
}
