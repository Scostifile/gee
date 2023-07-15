package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  // 完整的url，e.g. /p/:lang
	part     string  // 路由中的一部分	e.g. :lang
	children []*node // 子节点
	isWild   bool    // 是否模糊匹配 e.g. part中含有":filename"或"*filename"时为true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}",
		n.pattern, n.part, n.isWild)
}

// parts存放的是经pattern切分后的字段
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// parts存放的是经pattern切分后的字段
// bfs_search
func (n *node) search(parts []string, height int) *node {
	// 递归终止条件，找到末尾了或者通配符“*”【注意，对于通配符“:”也是需要找到末尾才能退出的】
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			// pattern为空字符串表示它不是一个完整的url，匹配失败
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	// bfs
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// 遍历所有完整的url，保存到列表中
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

// 第一个匹配成功的节点，用于insert
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于bfs_search
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
