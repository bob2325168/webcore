package middleware

import (
	"fmt"
	"webcore/framework"
)

func Foo1() framework.ControllerHandler {
	return func(ctx *framework.Context) error {
		fmt.Println("middleware pre test1")
		ctx.Next()
		fmt.Println("middleware post test1")
		return nil
	}
}

func Foo2() framework.ControllerHandler {
	return func(ctx *framework.Context) error {
		fmt.Println("middleware pre test2")
		ctx.Next()
		fmt.Println("middleware post test2")
		return nil
	}
}
