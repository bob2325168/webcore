package main

import (
	"net/http"
	"webcore"
	"webcore/framework"
	"webcore/middleware"
)

func main() {
	core := framework.NewCore()

	// 使用中间件
	core.Use(middleware.Recovery())
	core.Use(middleware.Cost())

	// 设置全局的middleware
	core.Use(middleware.Foo1())

	// 注册路由
	webcore.RegisterRouter(core)
	serve := &http.Server{ // nolint:gosec
		Addr:    ":8888",
		Handler: core,
	}
	serve.ListenAndServe()
}
