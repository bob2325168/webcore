package webcore

import (
	"webcore/framework"
	"webcore/middleware"
)

func RegisterRouter(core *framework.Core) {
	// HTTP方法+静态路由匹配
	//core.Get("/user/login", UserLoginControllerHandler)
	core.Get("/user/login", middleware.Foo1(), UserLoginControllerHandler)

	// 批量通用前缀匹配
	subjectAPI := core.Group("/subject")
	{
		subjectAPI.Get("/list/all", SubjectListControllerHandler)
		subjectAPI.Get("/:id", middleware.Foo2(), SubjectGetControllerHandler)
		subjectAPI.Put("/:id", SubjectUpdateControllerHandler)
		subjectAPI.Delete("/:id", SubjectDelControllerHandler)

		// 组内调用
		subjectInnerAPI := subjectAPI.Group("/info")
		{
			subjectInnerAPI.Get("/name", SubjectNameControllerHandler)
		}
	}
}
