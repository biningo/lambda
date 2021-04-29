package lambda

import (
	"log"
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
	log.Println("add ", fullPath)
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
	return cur
}

func (n *node) getHandlers(fullPath string) (HandlersChain, map[string]string) {
	fullPath = cleanPath(fullPath)
	node := n.getNode(fullPath)
	if node == nil {
		return HandlersChain{}, nil
	}
	patternPath := node.fullPath
	params := make(map[string]string)
	parts := strings.Split(patternPath, "/")[1:]
	pathParts := strings.Split(fullPath, "/")[1:]
	for index, part := range parts {
		if part[0] == ':' {
			params[part[1:]] = pathParts[index]
		}
	}
	return node.handlers, params
}
