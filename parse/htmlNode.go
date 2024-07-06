package parse

import "fmt"

type HTMLNode struct {
	tag      string     // HTML tag
	props    string     // Props to add to element TODO
	children []HTMLNode // HTML children
	value    string     // text
}

func (h HTMLNode) String() string {
	if h.children == nil {
		return fmt.Sprintf("LeafNode{\"<%s>\", props=%s, value=\"%s\"}", h.tag, h.props, h.value)
	} else {
		return fmt.Sprintf("ParentNode{\"<%s>\", props=%s, children=\"%s\"}", h.tag, h.props, h.children)
	}
}

func (h HTMLNode) ToHTML() string {
	if h.children == nil {
		textNodes := TextToTextNodes(h.value)
		value := TextNodesToHTML(textNodes)
		return fmt.Sprintf("<%s %s>%s</%s>", h.tag, h.props, value, h.tag)
	}

	var childHTML string
	for _, child := range h.children {
		childHTML += child.ToHTML()
	}
	return fmt.Sprintf("<%s %s>%s</%s>", h.tag, h.props, childHTML, h.tag)
}

func NewLeafNode(tag string, props string, value string) HTMLNode {
	return HTMLNode{
		tag:      tag,
		props:    props,
		children: nil,
		value:    value,
	}
}

func NewParentNode(tag string, props string, children []HTMLNode) HTMLNode {
	return HTMLNode{
		tag:      tag,
		props:    props,
		children: children,
		value:    "",
	}
}

