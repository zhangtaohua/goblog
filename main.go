package main

import (
	"net/http"

	"github.com/zhangtaohua/goblog/app/http/middlewares"
	"github.com/zhangtaohua/goblog/bootstrap"
	"github.com/zhangtaohua/goblog/config"
	c "github.com/zhangtaohua/goblog/pkg/config"
	"github.com/zhangtaohua/goblog/pkg/logger"
)

func init() {
	config.Initialize()
}

func main() {
	bootstrap.SetupDB()

	router := bootstrap.SetupRoute()

	err := http.ListenAndServe(":"+c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
	logger.LogError(err)
}

// Content-Type
// text/html
// text/plain
// text/css
// text/javascript
// application/json
// application/xml
// image/png
// go clean =modcache
// go clean -modcache
// go mod download
// go mod download
// go mod init
// go mod download
// go mod tidy
// go mod graph
// go mod edit
// go mod vendor
// go mod verify
// go mod verify
// go mod why
// GO111MODULE
// GOSUMDB
// GONOPROXY
// GOPROXY GONOPROXY GONOSUMDB GOPRIVATE GOPRIVATE GOPRIVATE GOENV
