package model

import (
	"net/http"
	"encoding/json"
)

type ResponseObj struct {
	Success bool `json:"success"`
	Error interface{} `json:"error,omitempty"`
	Body  interface{} `json:"body,omitempty"`
}

// Function write response object attributes as JSON to http response
func WriteResponse (success bool, errorMsg interface{}, body Entity, w http.ResponseWriter) {
	var entity interface{}
	if body != nil {
		entity = body.Entity()
	}

	var resp ResponseObj = ResponseObj{
		success,
		errorMsg,
		entity,
	}

	js, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json;utf-8")
	w.Write([]byte(js))
}

// Function write reponse object for array of entities
func WriteArrayResponse (success bool, errorMsg interface{}, body []Entity, w http.ResponseWriter) {
	var entities []interface{} = make([]interface{}, len(body))

	for index, entity := range body {
		entities[index] = entity.Entity()
	}

	var resp ResponseObj = ResponseObj{
		success,
		errorMsg,
		entities,
	}

	js, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json;utf-8")
	w.Write([]byte(js))
}
