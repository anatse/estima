package main

import (
	"log"
	"net/http"
	"os"
	"ru/sbt/estima/services"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"reflect"
	"unsafe"
	"encoding/json"
	"ru/sbt/estima/model"
)

func JwtHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer (func() {
			//if r := recover(); r != nil {
			//	//log.Println(r.(*errors.Error).ErrorStack())
			//	//fmt.Println("Recovered in Handler:", r)
			//	model.WriteResponse(false, fmt.Sprint(r), nil, w)
			//}
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
		rp.GetName()
		ris[i] = routeInfo{
			rp.GetName(),
			path,
		}
	}

	return ris
}

func main() {
	var us services.UserService
	var ps services.ProjectService
	var pcs services.ProcessService

	model.RegisterService("user", us)
	model.RegisterService("project", ps)
	model.RegisterService("process", pcs)

	r := mux.NewRouter()
	r.Handle("/get-token", services.GetTokenHandler).Methods("GET").Name("Login router (GET). Query parameters uname & upass")
	r.Handle("/login", services.Login).Methods("POST").Name("Login router (POST). Body: uname & upass")

	us.ConfigRoutes(r, JwtHandler)
	ps.ConfigRoutes(r, JwtHandler)
	pcs.ConfigRoutes(r, JwtHandler)

	// Function build router for get router information
	var routesInformation = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		js, _ := json.Marshal(getRoutes(r))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
