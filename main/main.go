package main

import (
	"net/http"
	"webcore"
	"webcore/framework"
)

func main() {
	core := framework.NewCore()
	// 注册路由
	webcore.RegisterRouter(core)
	serve := &http.Server{ //nolint:gosec
		Addr:    ":8888",
		Handler: core,
	}
	serve.ListenAndServe()
}
