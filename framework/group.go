package framework

// IGroup 前缀分组
type IGroup interface {
	Get(string, ControllerHandler)
	Put(string, ControllerHandler)
	Delete(string, ControllerHandler)
	Post(string, ControllerHandler)
}

type Group struct {
	core   *Core
	prefix string
}

func (g *Group) Get(url string, handler ControllerHandler) {
	url = g.prefix + url
	g.core.Get(url, handler)
}

func (g *Group) Put(url string, handler ControllerHandler) {
	url = g.prefix + url
	g.core.Put(url, handler)
}

func (g *Group) Delete(url string, handler ControllerHandler) {
	url = g.prefix + url
	g.core.Delete(url, handler)
}

func (g *Group) Post(url string, handler ControllerHandler) {
	url = g.prefix + url
	g.core.Post(url, handler)
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		prefix: prefix,
	}
}
