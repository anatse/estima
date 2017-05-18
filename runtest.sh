#!/usr/bin/env bash
export CONFIG_PATH=~/projects/estima/config.json
export DBJS_PATH=~/projects/estima/dbjs/

go test ./src/ru/... -v