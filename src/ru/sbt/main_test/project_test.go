package main_test

import (
	"testing"
	"ru/sbt/estima/model"
	"time"
	"net/http"
	"encoding/json"
)

type ProjectResponse struct {
	Success bool
	Error string
	Body model.Project
}

type ProjectArrayResponse struct {
	Success bool
	Error string
	Body []model.Project
}

type StageResponse struct {
	Success bool
	Error string
	Body model.Stage
}

type StageArrayResponse struct {
	Success bool
	Error string
	Body []model.Stage
}

const (
	PRJ_NUM = "000000"
	STAGE = "Stage1"
)

var prjKey string
var stageKey string

// Test project creation service
// URL: /project/create, method POST
func TestProjectCreate (t *testing.T) {
	var prj model.Project
	prj.Number = PRJ_NUM
	prj.Description = "Test Project"
	prj.Flag = "A"
	prj.Status = "TESTING"
	prj.StartDate = time.Now()
	prj.EndDate = prj.StartDate.Add(time.Hour * 24)

	response := callSecured(http.NewRequest("POST", "/project/create", CreateBody(prj)))
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "" {
		var resp ProjectResponse
		checkError(json.Unmarshal([]byte(body), &resp), t)

		if resp.Body.Key == "" {
			t.Errorf("Project key should not be empty")
			t.FailNow()
		} else {
			prjKey = resp.Body.Key
		}

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}

		if resp.Body.Number != PRJ_NUM {
			t.Errorf("Expected project number: " + PRJ_NUM + ". Got %v", resp.Body.Number)
			t.FailNow()
		}

		if resp.Body.Status != "TESTING" {
			t.Errorf("Expected project number: TESTING. Got %v", resp.Body.Status)
			t.FailNow()
		}

		if resp.Body.Description != "Test Project" {
			t.Errorf("Expected project number: 'Test Project'. Got %v", resp.Body.Description)
			t.FailNow()
		}

		if resp.Body.Flag != "A" {
			t.Errorf("Expected project number: 'Flag'. Got %v", resp.Body.Flag)
			t.FailNow()
		}

	}
}

// Test adding predefined user to project 000000
// URL: /project/{id}/user/add, method POST
func TestAddUserToProject (t *testing.T) {
	var userInfo struct {
		Name string `json:"name"`
		Role string `json:"role"`
	}

	userInfo.Name = USER_NAME
	userInfo.Role = "TEST_ROLE"

	body := CreateBody(userInfo)
	response := callSecured(http.NewRequest("POST", "/project/"+ prjKey + "/user/add", body))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp ProjectResponse
		checkError(json.Unmarshal([]byte(body), &resp), t)

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}
	}
}

// Test getting projects list for current user
// URL: /user/projects, method GET
func TestGetProjectsByUser (t *testing.T) {
	response := callSecured(http.NewRequest("GET", "/user/projects", nil))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp ProjectArrayResponse
		checkError(json.Unmarshal([]byte(body), &resp), t)

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}

		if len(resp.Body) == 0 {
			t.Errorf("Expected at least one project. Got %v", len(resp.Body))
			t.FailNow()
		}
	}
}

// Test getting list of all project's users
// URL: /project/{id}/user/list
func TestProjectUserList (t *testing.T) {
	response := callSecured(http.NewRequest("GET", "/project/" + prjKey + "/user/list", nil))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp UserArrayResponse
		checkError(json.Unmarshal([]byte(body), &resp), t)

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}

		if len(resp.Body) == 0 {
			t.Errorf("Expected at least one user. Got %v", len(resp.Body))
			t.FailNow()
		}
	}
}

// Test removing user from project
// URL: /project/{id}/user/remove, method DELETE
func TestRemoveUserFromProject (t *testing.T) {
	var userInfo struct {
		Name string `json:"name"`
		Role string `json:"role"`
	}

	userInfo.Name = USER_NAME
	userInfo.Role = "TEST_ROLE"

	body := CreateBody(userInfo)
	response := callSecured(http.NewRequest("DELETE", "/project/"+ prjKey + "/user/remove", body))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp ProjectResponse
		checkError(json.Unmarshal([]byte(body), &resp), t)

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}
	}
}

// Test user list check for empty lists, this test should be called after remove user test
func TestProjectEmptyUserList (t *testing.T) {
	response := callSecured(http.NewRequest("GET", "/project/" + prjKey + "/user/list", nil))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp UserArrayResponse
		checkError(json.Unmarshal([]byte(body), &resp), t)

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}

		if len(resp.Body) != 0 {
			t.Errorf("Expected no users in project. Got %v", len(resp.Body))
			t.FailNow()
		}
	}
}

// Test getting list of all projects
// URL: /project/list, method GET
func TestFindAllProjects (t *testing.T) {
	response := callSecured(http.NewRequest("GET", "/project/list", nil))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp ProjectArrayResponse
		checkError(json.Unmarshal([]byte(body), &resp), t)

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}

		if len(resp.Body) == 0 {
			t.Errorf("Expected at least one project. Got %v", len(resp.Body))
			t.FailNow()
		}
	}
}

// Testing project stages. This test for empty list of stages
// URL: /project/{id}/stage/list, method GET
func TestProjectStagesList (t *testing.T) {
	response := callSecured(http.NewRequest("GET", "/project/"+ prjKey + "/stage/list", nil))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp StageArrayResponse
		checkError(json.Unmarshal([]byte(body), &resp), t)

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}

		if len(resp.Body) != 0 {
			t.Errorf("Expected no stages in the prpoject. Got %v", len(resp.Body))
			t.FailNow()
		}
	}
}

// Test adding stage to project
// URL: /project/{id}/stage/add, method POST
func TestAddStageToProject (t *testing.T) {
	var stg model.Stage

	stg.Name = STAGE
	stg.Description = "First stage"
	stg.Status = "READY"
	stg.StartDate = time.Now()
	stg.EndDate = stg.StartDate.AddDate(1, 0, 0)

	body := CreateBody(stg)
	response := callSecured(http.NewRequest("POST", "/project/"+ prjKey + "/stage/add", body))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp struct {
			Success bool `json:"success"`
			Error string `json:"error"`
			EntityId string `json:"entityId"`
		}

		checkError(json.Unmarshal([]byte(body), &resp), t)

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}
	}
}

func TestGetStageByName (t *testing.T) {
	var stg model.Stage
	stg.Name = STAGE
	body := CreateBody(stg)
	response := callSecured(http.NewRequest("POST", "/project/"+ prjKey + "/stage/get", body))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp StageResponse
		checkError(json.Unmarshal([]byte(body), &resp), t)

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}

		if resp.Body.Name != STAGE {
			t.Errorf("Expected stage name=true. Got %v", STAGE, resp.Body.Name)
			t.FailNow()
		}

		stageKey = resp.Body.Key
	}
}

// Test removing stage from project
// URL: /project/{id}/stage/remove, method DELETE
func TestRemoveStageFromProject (t *testing.T) {
	var stg model.Stage

	stg.Name = STAGE
	stg.Description = "First stage"
	stg.Status = "READY"
	stg.StartDate = time.Now()
	stg.EndDate = stg.StartDate.AddDate(1, 0, 0)

	body := CreateBody(stg)
	response := callSecured(http.NewRequest("DELETE", "/project/"+ prjKey + "/stage/remove", body))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp StageResponse
		checkError(json.Unmarshal([]byte(body), &resp), t)

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}
	}
}

// Test getting project info by its number
// URL: /project/{id}, method GET
func TestGetProjectByNumber (t *testing.T) {
	response := callSecured(http.NewRequest("GET", "/project/" + prjKey, nil))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp ProjectResponse
		err := json.Unmarshal([]byte(body), &resp)
		checkError(err, t)

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}

		if resp.Body.Number != PRJ_NUM {
			t.Errorf("Expected project name = "+ PRJ_NUM + ". Got %v", resp.Body.Number)
			t.FailNow()
		}
	}
}

// Test setting project status
// URL: /project/{id}/status, method POST
func TestSetProjectStatus (t *testing.T) {
	var prj model.Project
	prj.Status = "CHANGED"

	response := callSecured(http.NewRequest("POST", "/project/" + prjKey + "/status", CreateBody(prj)))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp ProjectResponse
		checkError(json.Unmarshal([]byte(body), &resp), t)

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}

		if resp.Body.Number != PRJ_NUM {
			t.Errorf("Expected project name = "+ PRJ_NUM + ". Got %v", resp.Body.Number)
			t.FailNow()
		}

		if resp.Body.Status != "CHANGED" {
			t.Errorf("Expected project status = CHANGED. Got %v", resp.Body.Status)
			t.FailNow()
		}
	}
}