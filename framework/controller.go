package framework

// ControllerHandler 返回错误可以校验
type ControllerHandler func(ctx *Context) error
