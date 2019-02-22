#!/bin/sh

export GO111MODULE=on
DIR=$(cd $(dirname $0); pwd)
BIN_DIR=$(cd $(dirname $(dirname $0)); pwd)/bin

mkdir -p ${BIN_DIR}
go build -o ${BIN_DIR}/exec-command ${DIR}/exec-command/main.go
go build -o ${BIN_DIR}/http-prompt ${DIR}/http-prompt/main.go
go build -o ${BIN_DIR}/live-prefix ${DIR}/live-prefix/main.go
go build -o ${BIN_DIR}/simple-echo ${DIR}/simple-echo/main.go
go build -o ${BIN_DIR}/simple-echo-cjk-cyrillic ${DIR}/simple-echo/cjk-cyrillic/main.go
