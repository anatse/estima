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
	//"github.com/go-errors/errors"
)

//var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	user := model.GetUserFromRequest (w, r)
//	js, _ := json.Marshal(user)
//	w.Header().Set("Content-Type", "application/json;utf-8")
//	w.Write([]byte(js))
//})

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

func main() {
	r := mux.NewRouter()
	r.Handle("/get-token", services.GetTokenHandler).Methods("GET")

	// Add user routers
	var us services.UserService
	us.ConfigRoutes(r, JwtHandler)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./views")))

	//err := http.ListenAndServeTLS(":9443", "server.crt", "server.key", handlers.LoggingHandler(os.Stdout, r))
	err := http.ListenAndServe(":9080", handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
