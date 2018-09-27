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
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	reqJSON := `{}`
	req, err := http.NewRequest("POST", "/organizations", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != 400 {
		t.Errorf("Expect status to be 400, but got %+v", rr.Code)
	}

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
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	reqJSON := `{"name": "test"}`
	req, err := http.NewRequest("POST", "/organizations", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Errorf("Expect status to be 200, but got %+v", rr.Code)
	}

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

func TestWhenUpdateOrganizationWithoutNameShouldReturnError(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	reqJSON := `{}`
	req, err := http.NewRequest("PUT", "/organizations", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != 400 {
		t.Errorf("Expect status to be 400, but got %+v", rr.Code)
	}

	var res testutil.Error
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Error != "NAME_EMPTY" {
		t.Errorf("Expect err to be NAME_EMPTY, but got %+v", res.Error)
	}
}

func TestWhenUpdateOrganizationWhenNotFoundShouldReturnError(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	reqJSON := `{"id":"5e9707b1-0000-0000-0000-02d2cef27bd9","name":"test"}`
	req, err := http.NewRequest("PUT", "/organizations", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != 404 {
		t.Errorf("Expect status to be 404, but got %+v", rr.Code)
	}

	var res testutil.Error
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Error != "NOT_FOUND" {
		t.Errorf("Expect err to be NOT_FOUND, but got %+v", res.Error)
	}
}

func TestWhenUpdateOrganizationSuccessShouldReturnError(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	org := model.Organization{
		Name: "name",
	}
	if err := pgdb.Create(&org).Error; err != nil {
		panic(err)
	}

	data, err := json.Marshal(endpoints.UpdateOrganizationRequest{
		ID:   org.ID,
		Name: "test",
	})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PUT", "/organizations", bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Errorf("Expect status to be 200, but got %+v", rr.Code)
	}

	var res endpoints.UpdateOrganizationResponse
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.ID != org.ID {
		t.Errorf("Expect Name to be test, but got %+v", res.Name)
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
	if o.ID != org.ID {
		t.Errorf("Expect Name to be test, but got %+v", res.Name)
	}
	if o.Name != "test" {
		t.Errorf("Expect Name to be test, but got %+v", res.Name)
	}
}
