package model

import (
	"net/http"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"encoding/json"
	"reflect"
	"regexp"
	"fmt"
	"log"
	"github.com/gorilla/mux"
)

var router *mux.Router
func GetRouter() *mux.Router{
	if router == nil {
		router = mux.NewRouter()
	}

	return router
}

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

type AraError struct {
	Exception string `json:"exception"`
	Stacktrace []string `json:"stacktrace"`
	Error bool `json:"error"`
	Code int `json:"code"`
	ErrorNum int `json:"errorNum"`
	ErrorMessage string `json:"errorMessage"`
}

func GetAraError (err interface{}) interface{} {
	errFmt, _ := regexp.Compile(`(\{"exception":".*".+"error":.+"code":.+"errorNum":.+"errorMessage".+\})`)
	errorString := fmt.Sprint(err)

	errorStrings := errFmt.FindAllStringSubmatch(errorString, -1)
	log.Println(errorStrings)

	if len(errorStrings) > 0 && len(errorStrings[0]) > 0 {
		errorString = errorStrings[0][0]
		var ae AraError
		json.Unmarshal([]byte(errorString), &ae)

		log.Printf (`Parsed error:
	exception: %s
	code: %s
	errorNum: %s
	errorMessage: %s
	stackTrace: %s		`,
			ae.Exception,
			ae.Code,
			ae.ErrorNum,
			ae.ErrorMessage,
			ae.Stacktrace)

		return ae
	}

	return nil
}