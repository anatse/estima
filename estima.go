package main

import (
	"log"
	"net/http"
	"os"
	"ru/sbt/estima/services"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"ru/sbt/estima/model"
	"fmt"
	"reflect"
	"unsafe"
	"encoding/json"
)

func JwtHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer (func() {
			if r := recover(); r != nil {
				//log.Println(r.(*errors.Error).ErrorStack())
				//fmt.Println("Recovered in Handler:", r)
				model.WriteResponse(false, fmt.Sprint(r), nil, w)
			}
		})()

		// Let secure process the request. If it returns an error,
		// that indicates the request should not continue.
		err := services.JwtMiddleware.CheckJWT(w, r)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		h.ServeHTTP(w, r)
	})
}

// Router information struct
type routeInfo struct {
	Name string
	Path string
}

// Fill router information
func getRoutes (router *mux.Router) []routeInfo {
	v := reflect.ValueOf(router)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	routes := v.FieldByName("routes")
	len := routes.Len()

	ris := make ([]routeInfo, len)

	for i:=0;i<len;i++ {
		route := routes.Index(i)
		rp := (*mux.Route)(unsafe.Pointer (route.Pointer()))

		path, _ := rp.GetPathTemplate()
		ris[i] = routeInfo{
			rp.GetName(),
			path,
		}
	}

	return ris
}

func main() {
	r := mux.NewRouter()
	r.Handle("/get-token", services.GetTokenHandler).Methods("GET").Name("Login router")

	// Add user routers
	var us services.UserService
	us.ConfigRoutes(r, JwtHandler)

	// Add project routers
	var ps services.ProjectService
	ps.ConfigRoutes(r, JwtHandler)

	// Function build router for get router information
	var routesInformation = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		js, _ := json.Marshal(getRoutes(r))
		w.Header().Set("Content-Type", "application/json;utf-8")
		w.Write([]byte(js))
	})

	r.Handle("/ri", routesInformation)

	// Add static router. Should be last in routes list
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./views")))

	//err := http.ListenAndServeTLS(":9443", "server.crt", "server.key", handlers.LoggingHandler(os.Stdout, r))
	err := http.ListenAndServe(":9080", handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
