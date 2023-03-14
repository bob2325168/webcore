package middleware

import (
	"context"
	"fmt"
	"time"
	"webcore/framework"
)

func Timeout(d time.Duration) framework.ControllerHandler {

	// 使用函数回调
	return func(ctx *framework.Context) error {
		//处理panic的消息通知
		panicCh := make(chan any, 1)
		//处理完成的消息通知
		finishCh := make(chan struct{}, 1)

		//baseContext是继承request的Context
		durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), d)
		defer cancel()

		//处理panic
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicCh <- p
				}
			}()

			// 使用Next执行具体的业务逻辑
			ctx.Next()

			finishCh <- struct{}{}
		}()
		// 监听超时，异常以及结束事件
		select {
		// Panic事件
		case <-panicCh:
			ctx.SetStatus(500).JSON("time out")
		// 结束事件
		case <-finishCh:
			fmt.Println("finish")
		// 超时事件
		case <-durationCtx.Done():
			ctx.SetTimeout()
			ctx.SetStatus(500).JSON("time out")
		}
		return nil
	}
}
