#!/usr/bin/env bash
export GOPATH=~/projects/estima/

# Запускаем БД.
/usr/local/opt/arangodb/sbin/arangod --daemon --pid-file /usr/local/opt/arangodb/arangod.pid &

# Запускаем Estimator.
go run estima.go &

# Стенд разработки UI.
cd ./src-ui
npm run start
