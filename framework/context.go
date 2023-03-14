package framework

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter

	// 是否超时
	hasTimeout  bool
	ctx         context.Context
	writerMutex *sync.Mutex

	// 当前请求的handler链条
	handlers []ControllerHandler
	index    int

	// url路由匹配参数
	params map[string]string
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writerMutex:    &sync.Mutex{},
		index:          -1,
	}
}

// #region base controller

func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writerMutex
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

func (ctx *Context) SetTimeout() {
	ctx.hasTimeout = true
}

// #endregion

// #region context.Context

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key any) any {
	return ctx.BaseContext().Value(key)
}

// SetHandlers 为context设置handlers
func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlers = handlers
}

// SetParams 设置参数
func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}

/**
Next 会在2个地方调用
1. 请求处理的入口，ServeHTTP
2. 每个中间件的逻辑代码中，用于调用下一个中间件
*/
// Next 调用context的下一个函数
func (ctx *Context) Next() error {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		if err := ctx.handlers[ctx.index](ctx); err != nil {
			return err
		}
	}
	return nil
}
