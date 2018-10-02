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

func TestWhenCreateGroupWithoutNameShouldReturnError(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	reqJSON := `{}`
	req, err := http.NewRequest("POST", "/groups", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
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
	err = pgdb.Model(&model.Group{}).
		Count(&count).
		Error
	if err != nil {
		panic(err)
	}
	if count != 0 {
		t.Errorf("Expect count to be 0, but got %+v", count)
	}
}

func TestWhenCreateGroupWithOrgNotFound(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	reqJSON := `{"name": "group", "organization_id": "5e9707b1-0000-0000-0000-02d2cef27bd9"}`
	req, err := http.NewRequest("POST", "/groups", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expect status to be 404, but got %+v", rr.Code)
	}

	var res testutil.Error
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Error != "ORGANIZATION_NOT_FOUND" {
		t.Errorf("Expect err to be ORGANIZATION_NOT_FOUND, but got %+v", res.Error)
	}

	var count int
	err = pgdb.Model(&model.Group{}).
		Count(&count).
		Error
	if err != nil {
		panic(err)
	}
	if count != 0 {
		t.Errorf("Expect count to be 0, but got %+v", count)
	}
}

func TestWhenCreateGroupSuccess(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	org := model.Organization{Name: "org", Metadata: make(model.Metadata)}
	if err := pgdb.Create(&org).Error; err != nil {
		panic(err)
	}
	reqJSON := `{"name": "test", "metadata": {"t1": "1", "t2": "2"}, "organization_id": "` + org.ID.String() + `"}`
	req, err := http.NewRequest("POST", "/groups", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expect status to be 200, but got %+v", rr.Code)
	}

	var res endpoints.CreateGroupResponse
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Name != "test" {
		t.Errorf("Expect Name to be test, but got %+v", res.Name)
	}
	if res.Status != 1 {
		t.Errorf("Expect Status to be 1, but got %+v", res.Status)
	}
	if len(res.Metadata) != 2 {
		t.Errorf("Expect Metadata to have size=2, but got %+v", len(res.Metadata))
	}
	if res.Metadata["t1"] != "1" {
		t.Errorf("Expect Metadata to have t1=1, but got %+v", res.Metadata["t1"])
	}
	if res.Metadata["t2"] != "2" {
		t.Errorf("Expect Metadata to have t2=2, but got %+v", res.Metadata["t2"])
	}

	var count int
	err = pgdb.Model(&model.Group{}).
		Count(&count).
		Error
	if err != nil {
		panic(err)
	}
	if count != 1 {
		t.Errorf("Expect count to be 1, but got %+v", count)
	}

	g := model.Group{}
	pgdb.First(&g)
	if g.Name != "test" {
		t.Errorf("Expect Name to be test, but got %+v", res.Name)
	}
	if g.Status != 1 {
		t.Errorf("Expect Status to be 1, but got %+v", res.Status)
	}
	if len(g.Metadata) != 2 {
		t.Errorf("Expect Metadata to have size=2, but got %+v", len(g.Metadata))
	}
}

func TestWhenUpdateGroupWhenNotFoundShouldReturnError(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	reqJSON := `{"name":"test"}`
	req, err := http.NewRequest("PUT", "/groups/5e9707b1-0000-0000-0000-02d2cef27bd9", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
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

func TestWhenUpdateGroupSuccess(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	org := model.Organization{Name: "org", Metadata: make(model.Metadata)}
	if err := pgdb.Create(&org).Error; err != nil {
		panic(err)
	}

	gr := model.Group{
		Name:           "name",
		OrganizationID: org.ID,
		Metadata:       make(model.Metadata),
	}
	if err := pgdb.Create(&gr).Error; err != nil {
		panic(err)
	}

	reqJSON := `{"name":"test"}`
	req, err := http.NewRequest("PUT", "/groups/"+gr.ID.String(), bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expect status to be 200, but got %+v", rr.Code)
	}

	var res endpoints.UpdateGroupResponse
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.ID != gr.ID {
		t.Errorf("Expect ID to be %+v, but got %+v", gr.ID, res.Name)
	}
	if res.Status != 1 {
		t.Errorf("Expect Status to be 1, but got %+v", res.Status)
	}

	var count int
	err = pgdb.Model(&model.Group{}).
		Count(&count).
		Error
	if err != nil {
		panic(err)
	}
	if count != 1 {
		t.Errorf("Expect count to be 1, but got %+v", count)
	}

	g := model.Group{}
	pgdb.First(&g)
	if g.ID != gr.ID {
		t.Errorf("Expect Name to be test, but got %+v", res.Name)
	}
	if g.Name != "test" {
		t.Errorf("Expect Name to be test, but got %+v", res.Name)
	}
}
