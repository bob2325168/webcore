package framework

import (
	"log"
	"net/http"
	"strings"
)

// Core 框架核心结构
type Core struct {
	router map[string]*Tree
	// 设置中间件
	middlewares []ControllerHandler
}

func NewCore() *Core {
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()

	return &Core{
		router:      router,
		middlewares: nil,
	}
}

// 注册中间件
func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middlewares = append(c.middlewares, middlewares...)
}

// Group 初始化Group
// 这里返回的是一个约定，IGroup是一个接口协议，好处就是不依赖具体的 Group 实现
func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(r, w)

	// 查找路由节点
	node := c.FindRouterNode(r)
	if node == nil {
		// 没有找到，打印日志
		ctx.SetStatus(404).JSON("not found")
		return
	}

	// 设置context的handlers
	ctx.SetHandlers(node.handlers)

	//设置路由参数
	params := node.parseParamsFromEndNode(r.URL.Path)
	ctx.SetParams(params)

	// 调用路由函数，如果返回err，返回500
	if err := ctx.Next(); err != nil {
		ctx.SetStatus(500).JSON("inner error")
		return
	}
}

func (c *Core) Get(uri string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["GET"].AddRouter(uri, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Post(uri string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["POST"].AddRouter(uri, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Put(uri string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["PUT"].AddRouter(uri, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Delete(uri string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["DELETE"].AddRouter(uri, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 查找路由，如果没有匹配到返回nil
func (c *Core) FindRouterNode(req *http.Request) *node {
	// method转化为大写，保证大小写不敏感
	uri := req.URL.Path
	upperMethod := strings.ToUpper(req.Method)
	// 查找第一层
	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.root.matchNode(uri)
	}
	return nil
}
