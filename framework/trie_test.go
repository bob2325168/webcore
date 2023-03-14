package framework

import (
	"testing"
)

func Test_filterChildNodes(t *testing.T) {
	root := &node{
		isLast:  false,
		segment: "",
		handler: func(ctx *Context) error {
			return nil
		},
		childs: []*node{
			{
				isLast:  true,
				segment: "FOO",
				handler: func(ctx *Context) error {
					return nil
				},
				childs: nil,
			},
			{
				isLast:  false,
				segment: ":id",
				handler: func(ctx *Context) error {
					return nil
				},
				childs: nil,
			},
		},
	}

	{
		nodes := root.filterChildNodes("FOO")
		if len(nodes) != 2 {
			t.Error("foo error")
		}
	}

	{
		nodes := root.filterChildNodes(":foo")
		if len(nodes) != 2 {
			t.Error(":foo error")
		}
	}
}

func Test_matchNode(t *testing.T) {
	root := &node{
		segment: "",
		handler: nil,
		isLast:  false,
		childs: []*node{
			{
				segment: "FOO",
				handler: nil,
				isLast:  true,
				childs: []*node{
					&node{
						segment: "BAR",
						handler: func(ctx *Context) error {
							panic("bingo")
						},
						childs: []*node{},
						isLast: true,
					},
				},
			},
			{
				segment: ":id",
				handler: nil,
				isLast:  true,
				childs:  []*node{},
			},
		},
	}

	{
		node := root.matchNode("foo/bar")
		if node == nil {
			t.Error("match normal node error")
		}
	}
	{
		node := root.matchNode("test")
		if node == nil {
			t.Error("match test")
		}
	}
}
