package model

import (
	"net/http"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"encoding/json"
	"reflect"
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

func checkIsPointer (entity interface{}) {
	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr {
		panic ("Entity should be passed by reference")
	}
}

func ReadJsonBodyAny (r *http.Request, entity interface{})(error) {
	checkIsPointer (entity)
	bodySize := r.ContentLength
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, bodySize))
	CheckErr (err)

	return json.Unmarshal(body, entity)
}

func ReadJsonBody (r *http.Request, entity Entity) {
	checkIsPointer (entity)
	bodySize := r.ContentLength
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, bodySize))
	CheckErr (err)

	err = json.Unmarshal(body, entity)
	CheckErr (err)
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
	WriteResponse(true, "Not implemented yet", nil, w)
}

func CheckErr (err error) {
	if err != nil {
		panic(err)
	}
}