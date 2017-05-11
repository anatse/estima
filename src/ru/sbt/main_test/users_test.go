package main_test

import (
	"testing"
	"net/http"
	"encoding/json"
	"ru/sbt/estima/model"
	"net/http/httptest"
	"bytes"
	"ru/sbt/estima/services"
)

const (
	USER_NAME = "wps8admin"		// Predefined user name
)

type UserResponse struct {
	Success bool
	Error string
	Body model.EstimaUser
}

type UserArrayResponse struct {
	Success bool
	Error string
	Body []model.EstimaUser
}

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
		t.FailNow()
	}
}

func checkError (err error, t *testing.T) {
	if err != nil {
		t.Errorf("Error occurred %s", err)
		t.FailNow()
	}
}

func checkUserName (response *httptest.ResponseRecorder, t *testing.T) {
	if body := response.Body.String(); body != "" {
		var resp UserResponse

		err := json.Unmarshal([]byte(body), &resp)
		checkError(err, t)

		if resp.Body.Name != USER_NAME {
			t.Errorf("Expected an wps8admin. Got %s", resp.Body.Name)
			t.FailNow()
		}
	}
}

// Function login in GET implementation. Using predefined user wps8admin with LDAP fake configuration
// URL: /get-token, method GET
func login () *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", "/get-token?uname=" + USER_NAME + "&upass=xxx", nil)
	response := executeRequest(req)
	return response
}

// Function login in POST implementation. Using predefined user wps8admin with LDAP fake configuration
// URL: /login, method POST
func loginPost (t *testing.T) *httptest.ResponseRecorder {
	var credential struct {
		Uname string `json:"uname"`
		Upass string `json:"upass"`
	}

	credential.Uname = USER_NAME
	credential.Upass = "xxx"

	credB, err := json.Marshal(credential)
	checkError(err, t)
	body := bytes.NewReader(credB)

	req, _ := http.NewRequest("POST", "/login", body)
	response := executeRequest(req)
	return response
}

// Function log into system using user without any privileges - guest
func loginGuest () *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", "/get-token?uname=guest&upass=xxx", nil)
	response := executeRequest(req)
	return response
}

// Function calls secured service using cookie from called login function
func callSecured (req *http.Request, err error) *httptest.ResponseRecorder {
	response := login()
	authCookie := response.Result().Cookies()[0]
	req.AddCookie(authCookie)
	return executeRequest(req)
}

// Function calls secured service using cookie from called loginGuest function
func callGuestSecured (req *http.Request, err error) *httptest.ResponseRecorder {
	response := loginGuest()
	authCookie := response.Result().Cookies()[0]
	req.AddCookie(authCookie)
	return executeRequest(req)
}

// Test index creation. This function should create unique index user.name
func TestCreateIndex(t *testing.T) {
	dao := services.NewUserDao()
	err := dao.CreateIndexes("users")
	checkError(err, t)
}

// Test login for predefined user wps8admin, used for LDAP protocol = fake
// TODO in production mode this test should be fixed using real username/password for LDAP
func TestLogin(t *testing.T) {
	response := login()
	checkResponseCode(t, http.StatusOK, response.Code)
	checkUserName(response, t)
}

// Test login with POST for predefined user wps8admin, used for LDAP protocol = fake
func TestLoginPost (t *testing.T) {
	response := loginPost(t)
	checkResponseCode(t, http.StatusOK, response.Code)
	checkUserName(response, t)
}

// Testing function get current user
// Checks response code and user name for user wps8admin
// Usign login function to log ito the system
// URL: /users/current, method GET
func TestGetCurrentUser (t *testing.T) {
	response := callSecured (http.NewRequest("GET", "/users/current", nil))
	checkResponseCode(t, http.StatusOK, response.Code)
	checkUserName(response, t)
}

// Testing service which retrieves list of all users
// URL: /users/list, method GET
func TestGetUserList (t *testing.T) {
	response := callSecured (http.NewRequest("GET", "/users/list", nil))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp UserArrayResponse

		err := json.Unmarshal([]byte(body), &resp)
		checkError(err, t)

		if len(resp.Body) == 0 {
			t.Errorf("Expected at least one user. Got %d", len(resp.Body))
			t.FailNow()
		}
	}
}

// Testing service which searches users by name
// Using predefine user name to search 8adm - should find al least one user wps8admin
// URL: /users/search, method GET
func TestSearchPredefinedUsers (t *testing.T) {
	response := callSecured (http.NewRequest("GET", "/users/search?name=8adm", nil))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp UserArrayResponse

		err := json.Unmarshal([]byte(body), &resp)
		checkError(err, t)

		if resp.Success != true {
			t.Errorf("Expected success=true. Got suceess=%v", resp.Success)
			t.FailNow()
		}

		if len(resp.Body) == 0 {
			t.Errorf("Expected at least one user. Got %d", len(resp.Body))
			t.FailNow()
		}

		if resp.Body[0].Name != USER_NAME {
			t.Errorf("Expected user name wps8admin. Got %s", resp.Body[0].Name)
		}
	}
}

// Test checks response for forbidden operation
// We are using user without any privileges to search users
// URL: /users/search, method GET
func TestSearchUsersByGuest (t *testing.T) {
	response := callGuestSecured (http.NewRequest("GET", "/users/search?name=8adm", nil))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp UserArrayResponse

		err := json.Unmarshal([]byte(body), &resp)
		checkError(err, t)

		if resp.Success != false {
			t.Errorf("Expected success=false. Got suceess=%v", resp.Success)
			t.FailNow()
		}

		if len(resp.Body) != 0 {
			t.Errorf("Expected at least one user. Got %d", len(resp.Body))
			t.FailNow()
		}

		if resp.Error != "Insufficient privilegies" {
			t.Errorf("Expected error . Got %s", resp.Error)
			t.FailNow()
		}
	}
}

// Test user creation or update
func TestCreateUser (t *testing.T) {
	var user model.EstimaUser
	user.Name = "testUser"
	user.DisplayName = "User for Test purpose only"
	user.Roles = []string{"TEST_ROLE"}

	body := CreateBody (user)

	response := callSecured (http.NewRequest("POST", "/users/create", body))
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "" {
		var resp UserResponse

		err := json.Unmarshal([]byte(body), &resp)
		checkError(err, t)

		if resp.Success != true {
			t.Errorf("Expected success=false. Got suceess=%v, error=%v", resp.Success, resp.Error)
			t.FailNow()
		}

		if resp.Body.Name != "testUser" {
			t.Errorf("Expected testUser. Got %v", resp.Body.Name)
			t.FailNow()
		}

		if len(resp.Body.Roles) != 1 {
			t.Errorf("Expected number of roles 1. Got %v", len(resp.Body.Roles))
			t.FailNow()
		}

		if resp.Body.Roles[0] != "TEST_ROLE" {
			t.Errorf("Expected roles TEST_ROLE. Got %v", resp.Body.Roles[0])
			t.FailNow()
		}
	}
}
