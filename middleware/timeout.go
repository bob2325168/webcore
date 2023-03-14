package middleware

import (
	"context"
	"fmt"
	"time"
	"webcore/framework"
)

func Timeout(d time.Duration) framework.ControllerHandler {
	// 使用函数回调
	return func(c *framework.Context) error {
		//处理panic的消息通知
		panicCh := make(chan any, 1)
		//处理完成的消息通知
		finishCh := make(chan struct{}, 1)

		//baseContext是继承request的Context
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		//处理panic
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicCh <- p
				}
			}()
			// TODO
			// 使用Next执行具体的业务逻辑

			finishCh <- struct{}{}
		}()

		// 监听超时，异常以及结束事件
		select {
		// panic事件
		case <-panicCh:
			c.JSON(500, "time out")
		// 结束事件
		case <-finishCh:
			fmt.Println("complete")
		// 超时事件
		case <-durationCtx.Done():
			c.SetTimeout()
			c.JSON(500, "time out")
		}
		return nil
	}
}
