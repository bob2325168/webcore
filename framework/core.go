package framework

import (
	"net/http"
	"strings"
)

// Core 框架核心结构
type Core struct {
	router map[string]map[string]ControllerHandler
}

func NewCore() *Core {
	// 定义二级map
	getRouter := map[string]ControllerHandler{}
	postRouter := map[string]ControllerHandler{}
	putRouter := map[string]ControllerHandler{}
	deleteRouter := map[string]ControllerHandler{}

	// 将二级map写入一级map
	router := map[string]map[string]ControllerHandler{}
	router["GET"] = getRouter
	router["POST"] = postRouter
	router["PUT"] = putRouter
	router["DELETE"] = deleteRouter

	return &Core{
		router: router,
	}
}

// Group 初始化Group
// 这里返回的是一个约定，IGroup是一个接口协议，好处就是不依赖具体的 Group 实现
func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(r, w)
	// 查找路由
	handler := c.findRouter(r)
	if handler == nil {
		// 没有找到，打印日志
		ctx.JSON(404, "not found")
		return
	}
	// 调用路由函数，如果返回err，返回500
	if err := handler(ctx); err != nil {
		ctx.JSON(500, "inner error")
		return
	}
}

func (c *Core) Get(url string, handler ControllerHandler) {
	upperURL := strings.ToUpper(url)
	c.router["GET"][upperURL] = handler
}

func (c *Core) Post(url string, handler ControllerHandler) {
	upperURL := strings.ToUpper(url)
	c.router["POST"][upperURL] = handler
}

func (c *Core) Put(url string, handler ControllerHandler) {
	upperURL := strings.ToUpper(url)
	c.router["PUT"][upperURL] = handler
}

func (c *Core) Delete(url string, handler ControllerHandler) {
	upperURL := strings.ToUpper(url)
	c.router["DELETE"][upperURL] = handler
}

// 查找路由，如果没有匹配到返回nil
func (c *Core) findRouter(req *http.Request) ControllerHandler {
	// url和method都全部转化为大写，保证大小写不敏感
	upperUrl := strings.ToUpper(req.URL.Path)
	upperMethod := strings.ToUpper(req.Method)
	// 查找第一层
	if methodHandlers, ok := c.router[upperMethod]; ok {
		if controllerHandler, ok := methodHandlers[upperUrl]; ok {
			return controllerHandler
		}
	}
	return nil
}
