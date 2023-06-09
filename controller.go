package webcore

import (
	"context"
	"fmt"
	"time"
	"webcore/framework"
)

/**
如何使用自定义 Context 设置超时呢？结合前面分析的标准库思路，我们三步走完成：
1. 继承 request 的 Context，创建出一个设置超时时间的 Context；
2. 创建一个新的 Goroutine 来处理具体的业务逻辑；
3. 设计事件处理顺序，当前 Goroutine 监听超时时间 Contex 的 Done() 事件，和具体的业务处理结束事件，
哪个先到就先处理哪个。
*/

func FooControllerHandler(c *framework.Context) error {
	//处理panic的消息通知
	panicCh := make(chan any, 1)
	//处理完成的消息通知
	finishCh := make(chan struct{}, 1)

	//baseContext是继承request的Context
	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(1*time.Second))
	defer cancel()

	//处理panic
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicCh <- p
			}
		}()
		time.Sleep(10 * time.Second)
		c.SetOkStatus().JSON("ok")
		finishCh <- struct{}{}
	}()

	// 监听超时，异常以及结束事件
	select {
	// panic事件
	case <-panicCh:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.SetStatus(500).JSON("panic")
	// 结束事件
	case <-finishCh:
		fmt.Println("complete")
	// 超时事件
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.SetStatus(500).JSON("time out")
		c.SetTimeout()
	}
	return nil
}

func UserLoginControllerHandler(ctx *framework.Context) error {
	foo, _ := ctx.QueryString("foo", "def")
	// 等待10s结束执行
	time.Sleep(10 * time.Second)
	ctx.SetOkStatus().JSON("UserLoginController: " + foo)
	return nil
}

func SubjectListControllerHandler(ctx *framework.Context) error {
	ctx.SetOkStatus().JSON("SubjectListController")
	return nil
}

func SubjectDelControllerHandler(ctx *framework.Context) error {
	ctx.SetOkStatus().JSON("SubjectDelController")
	return nil
}

func SubjectGetControllerHandler(ctx *framework.Context) error {
	ctx.SetOkStatus().JSON("SubjectGetController")
	return nil
}

func SubjectUpdateControllerHandler(ctx *framework.Context) error {
	ctx.SetOkStatus().JSON("SubjectUpdateController")
	return nil
}

func SubjectNameControllerHandler(ctx *framework.Context) error {
	ctx.SetOkStatus().JSON("SubjectNameController")
	return nil
}
