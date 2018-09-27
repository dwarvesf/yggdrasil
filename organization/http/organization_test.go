package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dwarvesf/yggdrasil/organization/endpoints"
	"github.com/dwarvesf/yggdrasil/organization/model"
	"github.com/dwarvesf/yggdrasil/organization/util/testutil"
	"github.com/go-kit/kit/log"
)

func TestWhenCreateOrganizationWithoutNameShouldReturnError(t *testing.T) {
	reqJSON := `{}`
	req, err := http.NewRequest("POST", "/organizations", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var res testutil.Error
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Error != "NAME_EMPTY" {
		t.Errorf("Expect err to be NAME_EMPTY, but got %+v", res.Error)
	}

	var count int
	err = pgdb.Model(&model.Organization{}).
		Count(&count).
		Error
	if err != nil {
		panic(err)
	}
	if count != 0 {
		t.Errorf("Expect count to be 0, but got %+v", count)
	}
}

func TestWhenCreateOrganizationWithNameShouldReturnSuccess(t *testing.T) {
	reqJSON := `{"name": "test"}`
	req, err := http.NewRequest("POST", "/organizations", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var res endpoints.CreateOrganizationResponse
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Name != "test" {
		t.Errorf("Expect Name to be test, but got %+v", res.Name)
	}
	if res.Status != 1 {
		t.Errorf("Expect Status to be 1, but got %+v", res.Status)
	}

	var count int
	err = pgdb.Model(&model.Organization{}).
		Count(&count).
		Error
	if err != nil {
		panic(err)
	}
	if count != 1 {
		t.Errorf("Expect count to be 1, but got %+v", count)
	}

	o := model.Organization{}
	pgdb.First(&o)
	if o.Name != "test" {
		t.Errorf("Expect Name to be test, but got %+v", res.Name)
	}
	if o.Status != 1 {
		t.Errorf("Expect Status to be 1, but got %+v", res.Status)
	}
}
