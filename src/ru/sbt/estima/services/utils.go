package services

import (
	"net/http"
	"io"
	"ru/sbt/estima/model"
	"io/ioutil"
	"net/url"
	"strconv"
)

func ReadJsonBody (w http.ResponseWriter, r *http.Request, entity model.Entity)(model.Entity) {
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
