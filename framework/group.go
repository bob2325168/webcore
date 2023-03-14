package framework

// IGroup 前缀分组
type IGroup interface {
	// 实现httpMethod方法
	Get(string, ...ControllerHandler)
	Put(string, ...ControllerHandler)
	Delete(string, ...ControllerHandler)
	Post(string, ...ControllerHandler)

	// 实现嵌套group
	Group(string) IGroup

	// 实现嵌套中间件
	Use(middlewares ...ControllerHandler)
}

type Group struct {
	core   *Core  //维持一个core结构
	prefix string //group的通用前缀
	parent *Group //指向上一个group

	middlewares []ControllerHandler //存放中间件
}

func (g *Group) Get(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Get(uri, allHandlers...)
}

func (g *Group) Put(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Put(uri, allHandlers...)
}

func (g *Group) Delete(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Delete(uri, allHandlers...)
}

func (g *Group) Post(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Post(uri, allHandlers...)
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:        core,
		prefix:      prefix,
		parent:      nil,
		middlewares: []ControllerHandler{},
	}
}

func (g *Group) Group(uri string) IGroup {
	cgroup := NewGroup(g.core, uri)
	cgroup.parent = g
	return cgroup
}

// Use 注册中间件
func (g *Group) Use(middlewares ...ControllerHandler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// 获取当前group的绝对路径
func (g *Group) getAbsolutePrefix() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.getAbsolutePrefix() + g.prefix
}

// 获取某个group的middleware
func (g *Group) getMiddlewares() []ControllerHandler {
	if g.parent == nil {
		return g.middlewares
	}
	return append(g.parent.getMiddlewares(), g.middlewares...)
}
