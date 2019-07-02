package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dwarvesf/yggdrasil/services/organization/endpoints"
	"github.com/dwarvesf/yggdrasil/services/organization/model"
	"github.com/dwarvesf/yggdrasil/services/organization/util/testutil"
	uuid "github.com/satori/go.uuid"
)

func TestWhenCreateOrganizationWithoutNameShouldReturnError(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, true)

	reqJSON := `{}`
	req, err := http.NewRequest("POST", "/organizations", bytes.NewReader([]byte(reqJSON)))
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

func TestWhenCreateOrganizationWithNameAndEmptyMetadataShouldReturnSuccess(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, true)

	reqJSON := `{"name": "test"}`
	req, err := http.NewRequest("POST", "/organizations", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expect status to be 200, but got %+v", rr.Code)
	}

	var res endpoints.CreateOrganizationResponse
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Name != "test" {
		t.Errorf("Expect Name to be test, but got %+v", res.Name)
	}
	if len(res.Metadata) != 0 {
		t.Errorf("Expect Metadata to be {}, but got %+v", res.Metadata)
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
	if len(res.Metadata) != 0 {
		t.Errorf("Expect Metadata to be {}, but got %+v", res.Metadata)
	}
}

func TestWhenCreateOrganizationSuccess(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, true)

	reqJSON := `{"name": "test", "metadata": {"t1": "1", "t2": "2"}}`
	req, err := http.NewRequest("POST", "/organizations", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
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
	if len(o.Metadata) != 2 {
		t.Errorf("Expect Metadata to have size=2, but got %+v", len(o.Metadata))
	}
}

func TestWhenUpdateOrganizationWithInvalidID(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, true)

	reqJSON := `{"name": "test"}`
	req, err := http.NewRequest("PUT", "/organizations/5e9707b1", bytes.NewReader([]byte(reqJSON)))
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
	if res.Error != "INVALID_ID" {
		t.Errorf("Expect err to be INVALID_ID, but got %+v", res.Error)
	}
}

func TestWhenUpdateOrganizationWithoutNameShouldReturnError(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, true)

	reqJSON := `{}`
	req, err := http.NewRequest("PUT", "/organizations/5e9707b1-0000-0000-0000-02d2cef27bd9", bytes.NewReader([]byte(reqJSON)))
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
}

func TestWhenUpdateOrganizationWhenNotFoundShouldReturnError(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, true)

	reqJSON := `{"name":"test"}`
	req, err := http.NewRequest("PUT", "/organizations/5e9707b1-0000-0000-0000-02d2cef27bd9", bytes.NewReader([]byte(reqJSON)))
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

func TestWhenUpdateOrganizationSuccess(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, true)

	org := model.Organization{
		Name:     "name",
		Metadata: make(model.Metadata),
	}
	if err := pgdb.Create(&org).Error; err != nil {
		panic(err)
	}

	reqJSON := `{"name":"test"}`
	req, err := http.NewRequest("PUT", "/organizations/"+org.ID.String(), bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expect status to be 200, but got %+v", rr.Code)
	}

	var res endpoints.UpdateOrganizationResponse
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.ID != org.ID {
		t.Errorf("Expect ID to be %+v, but got %+v", org.ID, res.Name)
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
		t.Errorf("Expect %+v to be test, but got %+v", org.ID, o.ID)
	}
	if o.Name != "test" {
		t.Errorf("Expect Name to be test, but got %+v", res.Name)
	}
}

func TestArchiveOrganizationNotFound(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, true)

	req, err := http.NewRequest("POST", "/organizations/5e9707b1-0000-0000-0000-02d2cef27bd9/archive", bytes.NewReader([]byte("{}")))
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

func TestArchiveOrganizationSuccess(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, true)

	org := model.Organization{
		Name:     "name",
		Metadata: make(model.Metadata),
	}
	if err := pgdb.Create(&org).Error; err != nil {
		panic(err)
	}

	gr1 := model.Group{
		Name:           "name1",
		OrganizationID: org.ID,
		Metadata:       make(model.Metadata),
	}
	if err := pgdb.Create(&gr1).Error; err != nil {
		panic(err)
	}

	gr2 := model.Group{
		Name:           "name2",
		OrganizationID: org.ID,
		Metadata:       make(model.Metadata),
	}
	if err := pgdb.Create(&gr2).Error; err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "/organizations/"+org.ID.String()+"/archive", bytes.NewReader([]byte("{}")))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expect status to be 200, but got %+v", rr.Code)
	}

	var res endpoints.ArchiveOrganizationResponse
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.ID != org.ID {
		t.Errorf("Expect ID to be %+v, but got %+v", org.ID, res.Name)
	}
	if res.Status != 2 {
		t.Errorf("Expect Status to be 2, but got %+v", res.Status)
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
	pgdb.Where("id = ?", org.ID).First(&o)
	if o.Status != 2 {
		t.Errorf("Expect status to be 2, but got %+v", res.Status)
	}

	g := model.Group{}
	pgdb.Where("id = ?", gr1.ID).First(&g)
	if g.Status != 2 {
		t.Errorf("Expect status to be 2, but got %+v", g.Status)
	}
	pgdb.Where("id = ?", gr2.ID).First(&g)
	if g.Status != 2 {
		t.Errorf("Expect status to be 2, but got %+v", g.Status)
	}
}

func TestLeaveOrgSuccess(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, true)

	org := model.Organization{Name: "org", Metadata: make(model.Metadata)}
	if err := pgdb.Create(&org).Error; err != nil {
		panic(err)
	}

	userID := "5e9707b1-0000-0000-0000-02d2cef27bd9"
	rawID, err := uuid.FromString(userID)
	if err != nil {
		panic(err)
	}

	o := model.UserOrganizations{
		UserID:         rawID,
		OrganizationID: org.ID,
	}
	if err := pgdb.Create(&o).Error; err != nil {
		panic(err)
	}

	reqJSON := `{"user_id":"` + userID + `"}`
	req, err := http.NewRequest("POST", "/organizations/"+org.ID.String()+"/leave", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expect status to be 200, but got %+v", rr.Code)
	}

	var res endpoints.JoinOrganizationResponse
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Status != "success" {
		t.Errorf("Expect status to be success, but got %+v", res.Status)
	}

	var count int
	err = pgdb.Model(&model.UserOrganizations{}).
		Count(&count).
		Error
	if err != nil {
		panic(err)
	}
	if count != 1 {
		t.Errorf("Expect count to be 1, but got %+v", count)
	}

	var uo model.UserOrganizations
	err = pgdb.Model(&model.UserGroups{}).
		Where("user_id = ? AND organization_id = ?", userID, org.ID).
		First(&uo).
		Error
	if err != nil {
		panic(err)
	}
	if uo.LeftAt == nil {
		t.Errorf("Expect LeftAt not null, but got null")
	}
}

func TestLeaveOrgWhenHasNotJoined(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, true)

	org := model.Organization{Name: "org", Metadata: make(model.Metadata)}
	if err := pgdb.Create(&org).Error; err != nil {
		panic(err)
	}

	userID := "5e9707b1-0000-0000-0000-02d2cef27bd9"
	reqJSON := `{"user_id":"` + userID + `"}`
	req, err := http.NewRequest("POST", "/organizations/"+org.ID.String()+"/leave", bytes.NewReader([]byte(reqJSON)))
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
	if res.Error != "HAS_NOT_JOINED" {
		t.Errorf("Expect err to be ALREADY_JOINED, but got %+v", res.Error)
	}
}

func TestJoinOrgSuccess(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, true)

	org := model.Organization{Name: "org", Metadata: make(model.Metadata)}
	if err := pgdb.Create(&org).Error; err != nil {
		panic(err)
	}

	userID := "5e9707b1-0000-0000-0000-02d2cef27bd9"
	reqJSON := `{"user_id":"` + userID + `"}`
	req, err := http.NewRequest("POST", "/organizations/"+org.ID.String()+"/join", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expect status to be 200, but got %+v", rr.Code)
	}

	var res endpoints.JoinOrganizationResponse
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Status != "success" {
		t.Errorf("Expect status to be success, but got %+v", res.Status)
	}

	var count int
	err = pgdb.Model(&model.UserOrganizations{}).
		Count(&count).
		Error
	if err != nil {
		panic(err)
	}
	if count != 1 {
		t.Errorf("Expect count to be 1, but got %+v", count)
	}

	var uo model.UserOrganizations
	err = pgdb.Model(&model.UserOrganizations{}).
		Where("user_id = ? AND organization_id = ?", userID, org.ID).
		First(&uo).
		Error
	if err != nil {
		panic(err)
	}
}

func TestJoinOrgWhenAlreadyJoined(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, true)

	org := model.Organization{Name: "org", Metadata: make(model.Metadata)}
	if err := pgdb.Create(&org).Error; err != nil {
		panic(err)
	}

	org2 := model.Organization{Name: "org2", Metadata: make(model.Metadata)}
	if err := pgdb.Create(&org2).Error; err != nil {
		panic(err)
	}

	userID := "5e9707b1-0000-0000-0000-02d2cef27bd9"
	rawID, err := uuid.FromString(userID)
	if err != nil {
		panic(err)
	}

	uo := model.UserOrganizations{
		UserID:         rawID,
		OrganizationID: org.ID,
	}
	if err := pgdb.Create(&uo).Error; err != nil {
		panic(err)
	}

	reqJSON := `{"user_id":"` + userID + `"}`
	req, err := http.NewRequest("POST", "/organizations/"+org2.ID.String()+"/join", bytes.NewReader([]byte(reqJSON)))
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
	if res.Error != "ALREADY_JOINED" {
		t.Errorf("Expect err to be ALREADY_JOINED, but got %+v", res.Error)
	}
}

func TestJoinOrgWhenOrgNotActive(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, true)

	org := model.Organization{
		Name:     "org",
		Metadata: make(model.Metadata),
		Status:   model.OrganizationStatusInactive,
	}
	if err := pgdb.Create(&org).Error; err != nil {
		panic(err)
	}

	userID := "5e9707b1-0000-0000-0000-02d2cef27bd9"
	reqJSON := `{"user_id":"` + userID + `"}`
	req, err := http.NewRequest("POST", "/organizations/"+org.ID.String()+"/join", bytes.NewReader([]byte(reqJSON)))
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
	if res.Error != "ORG_NOT_ACTIVE" {
		t.Errorf("Expect err to be ORG_NOT_ACTIVE, but got %+v", res.Error)
	}
}

func TestJoinOrgWhenOrgNotFound(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, true)

	userID := "5e9707b1-0000-0000-0000-02d2cef27bd9"
	reqJSON := `{"user_id":"` + userID + `"}`
	req, err := http.NewRequest("POST", "/organizations/5e9707b1-1111-2222-3333-02d2cef27bd9/join", bytes.NewReader([]byte(reqJSON)))
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
