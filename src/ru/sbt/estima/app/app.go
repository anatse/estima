package app

import (
	"net/http"
	"ru/sbt/estima/model"
	"ru/sbt/estima/services"
	"github.com/gorilla/mux"
	"reflect"
	"unsafe"
	"encoding/json"
	"github.com/gorilla/handlers"
	"os"
	"log"
	"fmt"
)

func JwtHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer (func() {
			if r := recover(); r != nil {
				araErr := model.GetAraError (r)
				if araErr != nil {
					ae := araErr.(model.AraError)
					model.WriteResponse(false, fmt.Sprintf("%s", model.GetErrorText(ae)), nil, w)
				} else {
					model.WriteResponse(false, fmt.Sprint(r), nil, w)
				}
			}
		})()

		// Let secure process the request. If it returns an error,
		// that indicates the request should not continue.
		err := services.JwtMiddleware.CheckJWT(w, r)

		// If there was an error, do not continue.
		if err != nil {
			model.WriteResponse(false, fmt.Sprintf("Forbidden: %v", err), nil, w)
			w.WriteHeader(http.StatusOK)
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

func PrepareRoute () *mux.Router {
	var us services.UserService
	var ps services.ProjectService
	var pcs services.ProcessService
	var fs services.FeatureService

	model.RegisterService("user", us)
	model.RegisterService("project", ps)
	model.RegisterService("process", pcs)
	model.RegisterService("feature", fs)

	r := model.GetRouter()
	r.Handle("/api/v.0.0.1/get-token", services.GetTokenHandler).Methods("GET").Name("Login router (GET). Query parameters uname & upass")
	r.Handle("/api/v.0.0.1/login", services.Login).Methods("POST").Name("Login router (POST). Body: uname & upass")

	us.ConfigRoutes(r, JwtHandler)
	ps.ConfigRoutes(r, JwtHandler)
	pcs.ConfigRoutes(r, JwtHandler)
	fs.ConfigRoutes(r, JwtHandler)

	// Function build router for get router information
	var routesInformation = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		js, _ := json.Marshal(getRoutes(r))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(js))
	})

	r.Handle("/api/v.0.0.1/ri", routesInformation)

	// Add static router. Should be last in routes list
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./views")))
	return r
}

func AppRun () {
	// Init LDAP
	model.InitLdapPool(2)

	r := PrepareRoute()
	//err := http.ListenAndServeTLS(":9443", "server.crt", "server.key", handlers.LoggingHandler(os.Stdout, r))
	err := http.ListenAndServe(":9080", handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	model.FinishLdapPool()
}

