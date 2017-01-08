#!/bin/bash

# Install Go dependencies
go get -v -u github.com/jteeuwen/go-bindata/...
go get -v -u github.com/elazarl/go-bindata-assetfs/...
go get -v -u github.com/gorilla/mux
go get -v -u github.com/Sirupsen/logrus
go get -v -u github.com/rs/cors
go get -v -u github.com/spf13/viper
go get -v -u github.com/stianeikeland/go-rpio

# Build static content
pushd app/
[[ ! -d bower_components ]] && npm install && bower install
polymer build && rsync -av --delete build/unbundled/ ../static/
popd
go-bindata-assetfs static/...

GOOS=linux GOARCH=amd64 go build
