package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
	"ru/sbt/estima/model"
	"encoding/json"
	"log"
	"ru/sbt/estima/app"
)

// https://adazzle.github.io/react-data-grid/examples.html#/built-in-editors
// http://react-redux-grid.herokuapp.com/CustomLoader
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	model.GetRouter().ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func login () *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", "/get-token?uname=wps8admin&upass=xxx", nil)
	response := executeRequest(req)
	return response
}

func callSecured (req *http.Request, err error) *httptest.ResponseRecorder {
	response := login()
	authCookie := response.Result().Cookies()[0]
	req.AddCookie(authCookie)
	return executeRequest(req)
}

func checkUserName (response *httptest.ResponseRecorder, t *testing.T) {
	if body := response.Body.String(); body != "" {
		var resp struct {
			Success bool
			Body model.EstimaUser
		}

		err := json.Unmarshal([]byte(body), &resp)
		if err != nil {
			t.Errorf("Expected an empty array. Got %s", err)
		}

		if resp.Body.Name != `wps8admin` {
			t.Errorf("Expected an wps8admin. Got %s", resp.Body.Name)
		}
	}
}

func TestLogin(t *testing.T) {
	response := login()
	checkResponseCode(t, http.StatusOK, response.Code)
	checkUserName(response, t)
}

func TestGetCurrentUser (t *testing.T) {
	response := callSecured (http.NewRequest("GET", "/users/current", nil))
	checkResponseCode(t, http.StatusOK, response.Code)
	checkUserName(response, t)
}

func TestGetUserList (t *testing.T) {
	response := callSecured (http.NewRequest("GET", "/users/list", nil))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp struct {
			Success bool
			Body []model.EstimaUser
		}

		err := json.Unmarshal([]byte(body), &resp)
		if err != nil {
			t.Errorf("Expected an empty array. Got %s", err)
		}

		if len(resp.Body) == 0 {
			t.Errorf("Expected at least one user. Got %s", len(resp.Body))
		}
	}
}

func TestMain(m *testing.M) {
	log.Println("environment:" + os.Getenv("CONFIG_PATH"))
	app.PrepareRoute()
	code := m.Run()
	os.Exit(code)
}
