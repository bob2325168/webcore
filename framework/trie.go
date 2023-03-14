package framework

import (
	"errors"
	"fmt"
	"strings"
)

// Tree 代表树结构
type Tree struct {
	root *node // 代表根节点， 这个根节点是一个没有segment的空的根节点
}

// node 代表节点
type node struct {
	segment string            // uri中的字符串
	handler ControllerHandler // 控制器
	childs  []*node           // 子节点
	isLast  bool              // 该节点是不是一个独立的uri，是否自身就是一个终极节点
}

func newNode() *node {
	return &node{
		segment: "",
		handler: nil,
		childs:  []*node{},
		isLast:  false,
	}
}

func NewTree() *Tree {
	root := newNode()
	return &Tree{root: root}
}

/**
/book/list
/book/:id 冲突
/book/:id/name
/book/:student/age
/:user/name
/:user/name/:age 冲突
*/
// AddRouter 增加路由节点
func (t *Tree) AddRouter(uri string, handler ControllerHandler) error {
	n := t.root
	if n.matchNode(uri) != nil {
		return errors.New(fmt.Sprintf("trie: router exist %s", uri))
	}
	segments := strings.Split(uri, "/")
	for i, seg := range segments {
		if !isWildSegment(seg) {
			seg = strings.ToUpper(seg)
		}
		isLast := i == len(segments)-1
		//标记是否有合适的子节点
		var objNode *node
		childNodes := n.filterChildNodes(seg)
		if len(childNodes) > 0 {
			// 如果有segment相同的子节点，选择这个子节点
			for _, cnode := range childNodes {
				if cnode.segment == seg {
					objNode = cnode
					break
				}
			}
		}

		if objNode == nil {
			// 创建一个当前node的节点
			cnode := newNode()
			cnode.segment = seg
			if isLast {
				cnode.isLast = true
				cnode.handler = handler
			}
			n.childs = append(n.childs, cnode)
			objNode = cnode
		}

		n = objNode
	}
	return nil
}

// 匹配URI
func (t *Tree) FindHandler(uri string) ControllerHandler {
	matchedNode := t.root.matchNode(uri)
	if matchedNode == nil {
		return nil
	}
	return matchedNode.handler
}

// 判断一个segment是否是通用的segment，以 : 开头的，比如/user/:id
func isWildSegment(seg string) bool {
	return strings.HasPrefix(seg, ":")
}

// 判断路由是否已经在节点的所有子节点中存在
func (n *node) matchNode(uri string) *node {
	segments := strings.SplitN(uri, "/", 2)
	if len(segments) == 0 {
		return nil
	}
	// 第一部分用于匹配下一层子节点
	segment := segments[0]
	// 判断segment是不是通配符，以":"开头，例如/user/:id
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}
	// 匹配符合的下一层子节点
	nodes := n.filterChildNodes(segment)

	// 如果当前子节点没有符合的，说明这个uri之前不存在直接返回nil
	if nodes == nil || len(nodes) == 0 {
		return nil
	}
	// 如果只有一个segment，打上标记
	if len(segments) == 1 {
		for _, tn := range nodes {
			if tn.isLast {
				return tn
			}
		}
		// 都不是最后一个节点返回nil
		return nil
	}
	//如果有2个segment，递归每个子节点进行查找
	for _, tn := range nodes {
		matchedNode := tn.matchNode(segments[1])
		if matchedNode != nil {
			return matchedNode
		}
	}
	return nil
}

// 过滤所有的子节点
func (n *node) filterChildNodes(seg string) []*node {

	if len(n.childs) == 0 {
		return nil
	}
	// 如果segment是通配符，所以下一层子节点都满足需求
	if isWildSegment(seg) {
		return n.childs
	}

	nodes := make([]*node, 0, len(n.childs))
	// 过滤所有的下一层子节点
	for _, cnode := range n.childs {
		if isWildSegment(cnode.segment) {
			// 如果下一层子节点有通配符，满足需求
			nodes = append(nodes, cnode)
		} else if cnode.segment == seg {
			// 如果下一层子节点没有通配符，但是文本完全匹配，满足需求
			nodes = append(nodes, cnode)
		}
	}
	return nodes
}
