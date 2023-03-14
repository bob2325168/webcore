package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webcore"
	"webcore/framework"
	"webcore/middleware"
)

func main() {
	core := framework.NewCore()

	// 使用中间件
	core.Use(middleware.Recovery())
	core.Use(middleware.Cost())

	// 注册路由
	webcore.RegisterRouter(core)
	serve := &http.Server{ // nolint:gosec
		Addr:    ":8888",
		Handler: core,
	}
	// 开启2个goroutine，一个用于启动服务，一个用于监听信号（必须是 main 主goroutine）
	go func() {
		serve.ListenAndServe()
	}()

	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT SIGTERM SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 阻塞当前goroutine 等待信号
	<-quit

	// 调用serve.Shutdown 优雅关闭进程
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := serve.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
}
