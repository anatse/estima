package services

import (
	"net/http"
	"io"
	"ru/sbt/estima/model"
	"io/ioutil"
	"net/url"
	"strconv"
	"encoding/json"
)

var serviceMap map[string]interface{}
func RegisterService (name string, service interface{}) {
	if serviceMap == nil {
		serviceMap = make (map[string]interface{})
	}

	serviceMap[name] = service
}

func FindService (name string) interface{} {
	if serviceMap != nil {
		return serviceMap[name]
	} else {
		return nil
	}
}

func ReadJsonBodyAny (r *http.Request, entity interface{})(error) {
	bodySize := r.ContentLength
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, bodySize))
	if err != nil {
		panic (err)
	}

	return json.Unmarshal(body, entity)
}

func ReadJsonBody (r *http.Request, entity model.Entity)(model.Entity) {
	bodySize := r.ContentLength
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, bodySize))
	if err != nil {
		panic (err)
	}

	ret, err := entity.FromJson(body)
	if err != nil {
		panic(err)
	}

	return ret
}

func GetInt (values url.Values, name string,  def int) int {
	val := values.Get(name)
	if val != "" {
		iVal, _ := strconv.Atoi(val)
		return iVal
	}

	return def
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	model.WriteResponse(true, "Not implemented yet", nil, w)
}