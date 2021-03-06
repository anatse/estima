package main_test

import (
	"testing"
	"os"
	"log"
	"ru/sbt/estima/app"
	"encoding/json"
	"bytes"
	"io"
)

func CreateBody (value interface{}) io.Reader {
	v, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(v)
}

// Main testing function endpoint for whole testing
// to run tests
// 1. set environment variable CONFIG_PATH=<Full path to config,json file>
// 2. run command go test ./src/ru/... -v
func TestMain(m *testing.M) {
	log.Println("environment:" + os.Getenv("CONFIG_PATH"))
	app.PrepareRoute()
	code := m.Run()
	os.Exit(code)
}


