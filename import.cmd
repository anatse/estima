set GOPATH=E:\Dev\estima\

go get -u -v "github.com/dgrijalva/jwt-go"
go get -u -v "github.com/gorilla/context"
go get -u -v "github.com/gorilla/handlers"
go get -u -v "github.com/gorilla/mux"
go get -u -v "github.com/auth0/go-jwt-middleware"
go get -u -v "github.com/glycerine/zygomys/repl"
go get -u -v "gopkg.in/ldap.v2"
go get -u -v "github.com/go-errors/errors"

mkdir ./src/github.com/diegogub/
cd ./src/github.com/diegogub/
rd /s /q aranGO
git init
git clone https://github.com/anatse/aranGO.git
git clone https://github.com/diegogub/napping

npm install
npm run build
