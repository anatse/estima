#!/usr/bin/env bash
export GOPATH=~/projects/estima/

go get -u -v "github.com/dgrijalva/jwt-go"
go get -u -v "github.com/gorilla/context"
go get -u -v "github.com/gorilla/handlers"
go get -u -v "github.com/gorilla/mux"
go get -u -v "github.com/auth0/go-jwt-middleware"
go get -u -v "github.com/glycerine/zygomys/repl"
go get -u -v "gopkg.in/ldap.v2"
go get -u -v "github.com/go-errors/errors"
go get -u -v "github.com/bradfitz/gomemcache/memcache"
go get -u -v "github.com/gorilla/websocket"
go get -u -v "github.com/kardianos/service"

#
# AranGO should be clones from github because release too old therefore command go get "github.com/diegogub/aranGO" deprecated
#
mkdir ./src/github.com/diegogub/
cd ./src/github.com/diegogub/
rm -rf aranGO

git init
git clone https://github.com/anatse/aranGO.git
git clone https://github.com/diegogub/napping

cd ../../../src-ui
npm install
npm run build
