package lambda

import (
	"strings"
)

/**
*@Author lyer
*@Date 3/25/21 13:38
*@Describe
**/

type node struct {
	part     string
	fullPath string
	isWild   bool
	children []*node
	handlers HandlersChain
}

type Param struct {
	Key string
	Val string
}

//每个HTTP方法都有一棵树
type methodTree struct {
	method string
	root   *node
}

func (n *node) addRoute(fullPath string, handlers HandlersChain) {
	fullPath = cleanPath(fullPath)
	parts := strings.Split(fullPath, "/")[1:]
	cur := n
	for _, part := range parts {
		has := false
		for _, child := range cur.children {
			if part == child.part || child.isWild {
				has = true
				cur = child
				break
			}
		}
		if !has {
			newNode := &node{part: part}
			if part[0] == ':' {
				newNode.isWild = true
			}
			cur.children = append(cur.children, newNode)
			cur = newNode
		}
	}
	cur.fullPath = fullPath
	cur.handlers = handlers
}

func (n *node) getNode(fullPath string) *node {
	cur := n
	parts := strings.Split(fullPath, "/")[1:]
	for _, part := range parts {
		has := false
		for _, child := range cur.children {
			if part == child.part || child.isWild {
				has = true
				cur = child
				break
			}
		}
		if !has {
			return nil
		}
	}
	if cur.fullPath == fullPath {
		return cur
	}
	return nil
}

func (n *node) getHandlers(fullPath string) (HandlersChain, []Param) {
	fullPath = cleanPath(fullPath)
	node := n.getNode(fullPath)
	if node == nil {
		return HandlersChain{}, []Param{}
	}
	patternPath := node.fullPath
	params := []Param{}
	parts := strings.Split(patternPath, "/")[1:]
	pathParts := strings.Split(fullPath, "/")[1:]
	for index, part := range parts {
		if part[0] == ':' {
			params = append(params, Param{Key: part[1:], Val: pathParts[index]})
		}
	}
	return node.handlers, params
}
