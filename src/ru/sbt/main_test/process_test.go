package main_test

import (
	"testing"
	"ru/sbt/estima/model"
	"net/http"
	"encoding/json"
)

type ProcessResponse struct {
	Success bool
	Error string
	Body model.Process
}

type ProcessArrayResponse struct {
	Success bool
	Error string
	Body []model.Process
}

const (
	PRC_NAME = "Process1"
)

func TestPrepareProcess (t *testing.T) {
	TestProjectCreate(t)
	TestAddStageToProject(t)
	TestGetStageByName(t)
}

var prcKey string

// Test for process creation
// URL: /stage/{id}/process/create, method POST
func TestProcessCreate (t *testing.T) {
	var prc model.Process
	prc.Name = PRC_NAME
	prc.Description = "First Test Process"
	prc.Status = "TEST"

	response := callSecured(http.NewRequest("POST", "/stage/" + stageKey + "/process/create", CreateBody(prc)))
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "" {
		var resp ProcessResponse
		err := json.Unmarshal([]byte(body), &resp)
		checkError(err, t)

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}

		if resp.Body.Name != PRC_NAME {
			t.Errorf("Expected process name = %v. Got %v", PRC_NAME, resp.Body.Name)
			t.FailNow()
		}

		if resp.Body.Status != "TEST" {
			t.Errorf("Expected process status = TEST. Got %v", resp.Body.Status)
			t.FailNow()
		}

		if resp.Body.Key == "" {
			t.Errorf("Project key should not be empty")
			t.FailNow()
		} else {
			prcKey = resp.Body.Key
		}
	}
}

//
// URL: /process/{id}/remove, method DELETE
func TestProcessDelete (t *testing.T) {
	var prc model.Process
	prc.Name = PRC_NAME
	prc.Description = "First Test Process"
	prc.Status = "TEST"

	response := callSecured(http.NewRequest("DELETE", "/process/" + prcKey + "/remove", CreateBody(prc)))
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "" {
		var resp ProcessResponse
		err := json.Unmarshal([]byte(body), &resp)
		checkError(err, t)

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}

		if resp.Body.Name != PRC_NAME {
			t.Errorf("Expected process name = %v. Got %v", PRC_NAME, resp.Body.Name)
			t.FailNow()
		}
	}
}

func TestFinishProcess (t *testing.T) {
	TestRemoveStageFromProject(t)
}