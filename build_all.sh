#!/bin/bash

# Install Go dependencies
echo "Installing dependencies:"
echo -n "github.com/jteeuwen/go-bindata"
go get github.com/jteeuwen/go-bindata/...
echo -e " \033[0;32mOK\033[0m"
echo -n "github.com/elazarl/go-bindata-assetfs"
go get github.com/elazarl/go-bindata-assetfs/...
echo -e " \033[0;32mOK\033[0m"
echo -n "github.com/gorilla/mux"
go get github.com/gorilla/mux
echo -e " \033[0;32mOK\033[0m"
echo -n "github.com/Sirupsen/logrus"
go get github.com/Sirupsen/logrus
echo -e " \033[0;32mOK\033[0m"
echo -n "github.com/rs/cors"
go get github.com/rs/cors
echo -e " \033[0;32mOK\033[0m"
echo -n "github.com/spf13/viper"
go get github.com/spf13/viper
echo -e " \033[0;32mOK\033[0m"
echo -n "github.com/stianeikeland/go-rpio"
go get github.com/stianeikeland/go-rpio
echo -e " \033[0;32mOK\033[0m"

# Build static content
pushd app/
[[ ! -d bower_components ]] && npm install && bower install
polymer build && rsync -av --delete build/unbundled/ ../static/
popd
go-bindata-assetfs static/...

GOARM=7	GOARCH=arm go build
echo -e "\033[0;32mBUILD COMPLETED\033[0m"
