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
	"github.com/go-kit/kit/log"
	uuid "github.com/satori/go.uuid"
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
		t.Errorf("Expect ID to be %+v, but got %+v", gr.ID, g.ID)
	}
	if g.Name != "test" {
		t.Errorf("Expect Name to be test, but got %+v", res.Name)
	}
}

func TestWhenArchiveGroupSuccess(t *testing.T) {
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

	req, err := http.NewRequest("POST", "/groups/"+gr.ID.String()+"/archive", bytes.NewReader([]byte("{}")))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expect status to be 200, but got %+v", rr.Code)
	}

	var res endpoints.ArchiveGroupResponse
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.ID != gr.ID {
		t.Errorf("Expect ID to be %+v, but got %+v", gr.ID, res.Name)
	}
	if res.Status != 2 {
		t.Errorf("Expect Status to be 2, but got %+v", res.Status)
	}

	g := model.Group{}
	pgdb.First(&g)
	if g.ID != gr.ID {
		t.Errorf("Expect ID to be %+v, but got %+v", gr.ID, g.ID)
	}
	if g.Status != 2 {
		t.Errorf("Expect Status to be 2, but got %+v", res.Status)
	}
	if g.Name != "name" {
		t.Errorf("Expect Name to be name, but got %+v", res.Name)
	}
}

func TestLeaveGroupSuccess(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	org := model.Organization{Name: "org", Metadata: make(model.Metadata)}
	if err := pgdb.Create(&org).Error; err != nil {
		panic(err)
	}

	g := model.Group{
		Name:           "name",
		OrganizationID: org.ID,
		Metadata:       make(model.Metadata),
	}
	if err := pgdb.Create(&g).Error; err != nil {
		panic(err)
	}

	userID := "5e9707b1-0000-0000-0000-02d2cef27bd9"
	rawID, err := uuid.FromString(userID)
	if err != nil {
		panic(err)
	}

	u := model.UserGroups{
		UserID:  rawID,
		GroupID: g.ID,
	}
	if err := pgdb.Create(&u).Error; err != nil {
		panic(err)
	}

	reqJSON := `{"user_id":"` + userID + `"}`
	req, err := http.NewRequest("POST", "/groups/"+g.ID.String()+"/leave", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expect status to be 200, but got %+v", rr.Code)
	}

	var res endpoints.JoinGroupResponse
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Status != "success" {
		t.Errorf("Expect status to be success, but got %+v", res.Status)
	}

	var count int
	err = pgdb.Model(&model.UserGroups{}).
		Count(&count).
		Error
	if err != nil {
		panic(err)
	}
	if count != 1 {
		t.Errorf("Expect count to be 1, but got %+v", count)
	}

	var ug model.UserGroups
	err = pgdb.Model(&model.UserGroups{}).
		Where("user_id = ? AND group_id = ?", userID, g.ID).
		First(&ug).
		Error
	if err != nil {
		panic(err)
	}
	if ug.LeftAt == nil {
		t.Errorf("Expect LeftAt not null, but got null")
	}
}

func TestLeaveGroupWhenHasNotJoined(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	org := model.Organization{Name: "org", Metadata: make(model.Metadata)}
	if err := pgdb.Create(&org).Error; err != nil {
		panic(err)
	}

	g := model.Group{
		Name:           "name",
		OrganizationID: org.ID,
		Metadata:       make(model.Metadata),
	}
	if err := pgdb.Create(&g).Error; err != nil {
		panic(err)
	}

	userID := "5e9707b1-0000-0000-0000-02d2cef27bd9"
	reqJSON := `{"user_id":"` + userID + `"}`
	req, err := http.NewRequest("POST", "/groups/"+g.ID.String()+"/leave", bytes.NewReader([]byte(reqJSON)))
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

func TestJoinGroupSuccess(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	org := model.Organization{Name: "org", Metadata: make(model.Metadata)}
	if err := pgdb.Create(&org).Error; err != nil {
		panic(err)
	}

	g := model.Group{
		Name:           "name",
		OrganizationID: org.ID,
		Metadata:       make(model.Metadata),
	}
	if err := pgdb.Create(&g).Error; err != nil {
		panic(err)
	}

	userID := "5e9707b1-0000-0000-0000-02d2cef27bd9"
	reqJSON := `{"user_id":"` + userID + `"}`
	req, err := http.NewRequest("POST", "/groups/"+g.ID.String()+"/join", bytes.NewReader([]byte(reqJSON)))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expect status to be 200, but got %+v", rr.Code)
	}

	var res endpoints.JoinGroupResponse
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Status != "success" {
		t.Errorf("Expect status to be success, but got %+v", res.Status)
	}

	var count int
	err = pgdb.Model(&model.UserGroups{}).
		Count(&count).
		Error
	if err != nil {
		panic(err)
	}
	if count != 1 {
		t.Errorf("Expect count to be 1, but got %+v", count)
	}

	var ug model.UserGroups
	err = pgdb.Model(&model.UserGroups{}).
		Where("user_id = ? AND group_id = ?", userID, g.ID).
		First(&ug).
		Error
	if err != nil {
		panic(err)
	}
}

func TestJoinGroupWhenAlreadyJoined(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	org := model.Organization{Name: "org", Metadata: make(model.Metadata)}
	if err := pgdb.Create(&org).Error; err != nil {
		panic(err)
	}

	g := model.Group{
		Name:           "name",
		OrganizationID: org.ID,
		Metadata:       make(model.Metadata),
	}
	if err := pgdb.Create(&g).Error; err != nil {
		panic(err)
	}

	userID := "5e9707b1-0000-0000-0000-02d2cef27bd9"
	rawID, err := uuid.FromString(userID)
	if err != nil {
		panic(err)
	}

	u := model.UserGroups{
		UserID:  rawID,
		GroupID: g.ID,
	}
	if err := pgdb.Create(&u).Error; err != nil {
		panic(err)
	}

	reqJSON := `{"user_id":"` + userID + `"}`
	req, err := http.NewRequest("POST", "/groups/"+g.ID.String()+"/join", bytes.NewReader([]byte(reqJSON)))
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

func TestJoinGroupWhenGroupNotActive(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	org := model.Organization{Name: "org", Metadata: make(model.Metadata)}
	if err := pgdb.Create(&org).Error; err != nil {
		panic(err)
	}

	g := model.Group{
		Name:           "name",
		OrganizationID: org.ID,
		Status:         model.GroupStatusInactive,
		Metadata:       make(model.Metadata),
	}
	if err := pgdb.Create(&g).Error; err != nil {
		panic(err)
	}

	userID := "5e9707b1-0000-0000-0000-02d2cef27bd9"
	reqJSON := `{"user_id":"` + userID + `"}`
	req, err := http.NewRequest("POST", "/groups/"+g.ID.String()+"/join", bytes.NewReader([]byte(reqJSON)))
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
	if res.Error != "GROUP_NOT_ACTIVE" {
		t.Errorf("Expect err to be GROUP_NOT_ACTIVE, but got %+v", res.Error)
	}
}

func TestJoinGroupWhenGroupNotFound(t *testing.T) {
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	handler := NewHTTPHandler(pgdb, log.NewNopLogger(), true)

	userID := "5e9707b1-0000-0000-0000-02d2cef27bd9"
	reqJSON := `{"user_id":"` + userID + `"}`
	req, err := http.NewRequest("POST", "/groups/5e9707b1-1111-2222-3333-02d2cef27bd9/join", bytes.NewReader([]byte(reqJSON)))
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
