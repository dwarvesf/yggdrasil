package organization

import "net/http"

var (
	//ErrNameEmpty check Name not empty
	ErrNameEmpty = errNameEmpty{}
)

type errOrganization struct{}

func (errOrganization) Error() string {
	return "Organization invalid"
}

func (errOrganization) StatusCode() int {
	return http.StatusBadRequest
}

type errNameEmpty struct{}

func (errNameEmpty) Error() string {
	return "Name empty"
}

func (errNameEmpty) StatusCode() int {
	return http.StatusBadRequest
}
