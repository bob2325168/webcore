package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
	"webcore/err"
)

type Context struct {
	request *http.Request
	writer  http.ResponseWriter

	// 是否超时
	hasTimeout bool
	ctx        context.Context

	writerMutex *sync.Mutex
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:     r,
		writer:      w,
		ctx:         r.Context(),
		writerMutex: &sync.Mutex{},
	}
}

// #region base controller

func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writerMutex
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.writer
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

// #endregion

// #region form post

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request != nil {
		return ctx.request.PostForm
	}
	return map[string][]string{}
}
func (ctx *Context) FormInt(key string, def int) int {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			intVal, err := strconv.Atoi(vals[len-1])
			if err != nil {
				return def
			}
			return intVal
		}
	}
	return def
}

func (ctx *Context) FormString(key string, def string) string {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return def
}

func (ctx *Context) FormArray(key string, def []string) []string {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

// #endregion

// #region query url

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		return ctx.request.URL.Query()
	}
	return map[string][]string{}
}
func (ctx *Context) QueryInt(key string, def int) int {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			intVal, err := strconv.Atoi(vals[len-1])
			if err != nil {
				return def
			}
			return intVal
		}
	}
	return def
}

func (ctx *Context) QueryString(key string, def string) string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return def
}

func (ctx *Context) QueryArray(key string, def []string) []string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

// #endregion

// #region application/json

func (ctx *Context) BindJSON(obj any) error {
	if ctx.request != nil {
		body, err := io.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		// copy body
		ctx.request.Body = io.NopCloser(bytes.NewBuffer(body))
		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return err.ErrRequestEmpty
	}
	return nil
}

// #endregion

// #region response

// JSON 处理json response
func (ctx *Context) JSON(code int, obj any) error {
	// 超时返回
	if ctx.HasTimeout() {
		return nil
	}
	ctx.writer.Header().Set("Content-Type", "application/json")
	ctx.writer.WriteHeader(code)
	byt, err := json.Marshal(obj)
	if err != nil {
		ctx.writer.WriteHeader(500)
		return err
	}
	ctx.writer.Write(byt)
	return nil
}

func (ctx *Context) HTML(code int, obj any, template string) error {
	return nil
}

func (ctx *Context) Text(code int, obj any) error {
	return nil
}

// #endregion
