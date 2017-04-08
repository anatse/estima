package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"ru/sbt/estima/model"
	"ru/sbt/estima/services"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user")
	fmt.Println("products....")

	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	eUser := model.NewUser(
		claims["name"].(string),
		claims["mail"].(string),
		"",
		claims["displayName"].(string),
		claims["uid"].(string),
	)

	js, _ := json.Marshal(eUser)
	w.Header().Set("Content-Type", "application/json;utf-8")
	w.Write([]byte(js))
})

func main() {
	r := mux.NewRouter()
	r.Handle("/get-token", services.GetTokenHandler).Methods("GET")
	r.Handle("/products", services.JwtMiddleware.Handler(NotImplemented)).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./views")))

	err := http.ListenAndServeTLS(":9443", "server.crt", "server.key", handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
