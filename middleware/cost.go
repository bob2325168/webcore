package middleware

import (
	"log"
	"time"
	"webcore/framework"
)

// Cost 记录请求耗时
func Cost() framework.ControllerHandler {
	return func(ctx *framework.Context) error {

		start := time.Now()
		ctx.Next()
		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api uri: %v, cost: %v ", ctx.GetRequest().RequestURI, cost.Seconds())

		return nil
	}
}
