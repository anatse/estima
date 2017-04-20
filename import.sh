#!/usr/bin/env bash
export GOPATH=~/projects/estima/

go get "github.com/dgrijalva/jwt-go"
go get "github.com/gorilla/context"
go get "github.com/gorilla/handlers"
go get "github.com/gorilla/mux"
go get "github.com/auth0/go-jwt-middleware"
go get "github.com/glycerine/zygomys/repl"
go get "gopkg.in/ldap.v2"
go get "github.com/go-errors/errors"

#
# AranGO should be clones from github because release too old therefore command go get "github.com/diegogub/aranGO" deprecated
#
mkdir ./src/github.com/diegogub/
cd ./src/github.com/diegogub/
git clone https://github.com/victorspringer/aranGO.git