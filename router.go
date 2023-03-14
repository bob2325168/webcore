package webcore

import "webcore/framework"

func RegisterRouter(core *framework.Core) {
	// HTTP方法+静态路由匹配
	core.Get("/user/login", UserLoginControllerHandler)
	// 批量通用前缀匹配
	subjectApi := core.Group("/subject")
	{
		subjectApi.Get("/list/all", SubjectListControllerHandler)
		subjectApi.Get("/:id", SubjectGetControllerHandler)
		subjectApi.Put("/:id", SubjectUpdateControllerHandler)
		subjectApi.Delete("/:id", SubjectDelControllerHandler)
	}
}
