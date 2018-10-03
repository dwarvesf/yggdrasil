package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/organization/endpoints"
	"github.com/dwarvesf/yggdrasil/organization/model"
	"github.com/dwarvesf/yggdrasil/organization/util/testutil"
	"github.com/go-kit/kit/log"
)

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
