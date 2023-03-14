package middleware

import "webcore/framework"

// Recovery 捕获协程中的函数异常
func Recovery() framework.ControllerHandler {
	return func(ctx *framework.Context) error {
		defer func() {
			if err := recover(); err != nil {
				ctx.JSON(500)
			}
		}()

		// 使用next执行具体业务逻辑
		ctx.Next()

		return nil
	}

}
