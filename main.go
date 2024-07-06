package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	fmt.Println("Static Gopher")

	args := os.Args[1:]

	if len(args) != 3 {
		fmt.Println("[!] Invalid usage - exiting.")
		os.Exit(0)
	}

	sourceDir := args[0]
	destDir := args[1]
	templateDir := args[2]

	fmt.Println(sourceDir, destDir, templateDir)

	textNodes := TextToTextNodes("This has *italic* text.")
    _ = textNodes
}

type BlockType int

const (
	PARAGRAPH BlockType = iota
	HEADING
	CODE
	QUOTE
	UNORDERED_LIST
	ORDERED_LIST
)

func (bt BlockType) Name() string {
	if bt == PARAGRAPH {
		return "PARAGRAPH"
	}
	if bt == HEADING {
		return "HEADING"
	}
	if bt == CODE {
		return "CODE"
	}
	if bt == QUOTE {
		return "QUOTE"
	}
	if bt == UNORDERED_LIST {
		return "UNORDERED_LIST"
	}
	if bt == ORDERED_LIST {
		return "ORDERED_LIST"
	}
	return "INVALID"
}

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
		value := textNodesToHTML(textNodes)
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

type TextType int

const (
	TEXT TextType = iota
	BOLD
	ITALIC
	CODE_TEXT
	LINK
	IMAGE
)

func (tt TextType) Name() string {
	if tt == TEXT {
		return "TEXT"
	}
	if tt == BOLD {
		return "BOLD"
	}
	if tt == ITALIC {
		return "ITALIC"
	}
	if tt == CODE_TEXT {
		return "CODE_TEXT"
	}
	if tt == LINK {
		return "LINK"
	}
	if tt == IMAGE {
		return "IMAGE"
	}
	return "INVALID"
}

type TextNode struct {
	text     string
	textType TextType
	url      string
}

func (t TextNode) String() string {
	return fmt.Sprintf("TextNode{\"%s\", %s, url=\"%s\"}", t.text, t.textType.Name(), t.url)
}

func (t TextNode) Tag() string {
	if t.textType == TEXT {
		return ""
	}
	if t.textType == BOLD {
		return "b"
	}
	if t.textType == ITALIC {
		return "i"
	}
	if t.textType == CODE_TEXT {
		return "code"
	}
	if t.textType == LINK {
		return "a"
	}
	if t.textType == IMAGE {
		return "img"
	}
	return ""
}

func (t TextNode) ToHTML() string {
	switch t.textType.Name() {
	case "BOLD":
		return fmt.Sprintf("<b>%s</b>", t.text)
	case "ITALIC":
		return fmt.Sprintf("<i>%s</i>", t.text)
	case "CODE_TEXT":
		return fmt.Sprintf("<code>%s</code>", t.text)
	case "LINK":
		return fmt.Sprintf("<a href=\"%s\">%s</a>", t.url, t.text)
	case "IMAGE":
		return fmt.Sprintf("<img src=\"%s\" alt=\"%s\" />", t.url, t.text)
	default:
		return t.text
	}

	return fmt.Sprintf("")
}

func textNodesToHTML(nodes []TextNode) string {
	var html string
	for _, node := range nodes {
		html += node.ToHTML()
	}
	return html
}

func NewTextNode(text string, textType TextType, url string) TextNode {
	return TextNode{
		text:     text,
		textType: textType,
		url:      url,
	}
}

func TextToTextNodes(text string) []TextNode {
	node := NewTextNode(text, TEXT, "")

	nodes, err := splitNodesDelimiter([]TextNode{node}, "**")
	if err != nil {
		log.Panic("splitting on bold:", err)
	}
	
    nodes, err = splitNodesDelimiter(nodes, "*")
	if err != nil {
		log.Panic("splitting on italic:", err)
	}
	
    nodes, err = splitNodesDelimiter(nodes, "`")
	if err != nil {
		log.Panic("splitting on code:", err)
	}
	
    nodes, err = splitNodesImages(nodes)
	if err != nil {
		log.Panic("splitting on images:", err)
	}
	
    nodes, err = splitNodesLinks(nodes)
    if err != nil {
		log.Panic("splitting on links:", err)
	}

	return nodes
}

func splitNodesDelimiter(nodes []TextNode, delimiter string) ([]TextNode, error) {
	var (
		newNodes   []TextNode
		re         *regexp.Regexp
		targetType TextType
	)

	if delimiter == "*" {
		re = regexp.MustCompile(`(.*)\*(.+)\*(.*)`)
		targetType = ITALIC
	} else if delimiter == "**" {
		re = regexp.MustCompile(`(.*)\*\*(.+)\*\*(.*)`)
		targetType = BOLD
	} else if delimiter == "`" {
		// Using double quotes to allow '`' in pattern
		re = regexp.MustCompile("(.*)`(.+)`(.*)")
		targetType = CODE_TEXT
	} else {
		return nil, errors.New(fmt.Sprintf("invalid delimiter \"%s\"", delimiter))
	}

	for _, node := range nodes {
		newNodes = append(newNodes, splitter(re, targetType, node)...)
	}

	return newNodes, nil
}

func splitter(re *regexp.Regexp, targetType TextType, node TextNode) []TextNode {
	result := re.FindStringSubmatch(node.text)

	if result == nil {
		return []TextNode{node}
	}

	var splitNodes []TextNode

	pre := result[1]
	if pre != "" {
		tempNode := NewTextNode(pre, node.textType, node.url)
		splitNodes = append(splitNodes, splitter(re, targetType, tempNode)...)
	}

	targetText := result[2]
	splitNodes = append(splitNodes, NewTextNode(targetText, targetType, ""))

	post := result[3]
	if post != "" {
		splitNodes = append(splitNodes, NewTextNode(post, node.textType, node.url))
	}

	return splitNodes
}

func splitNodesLinks(nodes []TextNode) ([]TextNode, error) {
    var splitNodes []TextNode
    linkPattern := regexp.MustCompile(`(.*)\[([^\[\]\(\)]+)\]\(([^\[\]\(\)]+)\)(.*)`)

    for _, node := range nodes {
        result := linkPattern.FindStringSubmatch(node.text)

        if result == nil {
            splitNodes = append(splitNodes, node)
            continue
        }

        pre := result[1]
        if pre != "" {
            tempNode := NewTextNode(pre, node.textType, node.url)
            preNodes, _ := splitNodesLinks([]TextNode{tempNode})
            splitNodes = append(splitNodes, preNodes...)
        }

        text := result[2]
        url := result[3]
        splitNodes = append(splitNodes, NewTextNode(text, LINK, url))

        post := result[4]
        if post != "" {
            splitNodes = append(splitNodes, NewTextNode(post, node.textType, node.url))
        }
    }

	return splitNodes, nil
}

func splitNodesImages(nodes []TextNode) ([]TextNode, error) {
    var splitNodes []TextNode
    imgPattern := regexp.MustCompile(`(.*)\!\[([^\[\]\(\)]+)\]\(([^\[\]\(\)]+)\)(.*)`)

    for _, node := range nodes {
        result := imgPattern.FindStringSubmatch(node.text)

        if result == nil {
            splitNodes = append(splitNodes, node)
            continue
        }

        pre := result[1]
        if pre != "" {
            tempNode := NewTextNode(pre, node.textType, node.url)
            preNodes, _ := splitNodesImages([]TextNode{tempNode})
            splitNodes = append(splitNodes, preNodes...)
        }

        text := result[2]
        url := result[3]
        splitNodes = append(splitNodes, NewTextNode(text, IMAGE, url))

        post := result[4]
        if post != "" {
            splitNodes = append(splitNodes, NewTextNode(post, node.textType, node.url))
        }
    }

	return splitNodes, nil
}

func getBlockType(blockText string) BlockType {
	headingPattern := regexp.MustCompile(`#+ (.*)`)
	result := headingPattern.FindStringSubmatch(blockText)
	if result != nil {
		return HEADING
	}

	if strings.HasPrefix(blockText, "```") && strings.HasSuffix(blockText, "```") {
		return CODE
	}

	isQuote := true
	for _, line := range strings.Split(blockText, "\n") {
		if !strings.HasPrefix(line, "> ") {
			isQuote = false
			break
		}
	}
	if isQuote {
		return QUOTE
	}

	isUl := true
	for _, line := range strings.Split(blockText, "\n") {
		if !strings.HasPrefix(line, "* ") && !strings.HasPrefix(blockText, "- ") {
			isUl = false
			break
		}
	}
	if isUl {
		return UNORDERED_LIST
	}

	isOl := true
	index := 1
	for _, line := range strings.Split(blockText, "\n") {
		if !strings.HasPrefix(line, fmt.Sprintf("%d. ", index)) {
			isOl = false
			break
		}
		index++
	}
	if isOl {
		return ORDERED_LIST
	}

	return PARAGRAPH
}

func MarkdownToBlocks(md string) []string {
	var textBlocks []string
	var current []string

	for _, chunk := range strings.Split(md, "\n") {
		if chunk != "" {
			current = append(current, chunk)
		} else {
			if current != nil {
				textBlocks = append(textBlocks, strings.Join(current, "\n"))
				current = nil
			}
		}
	}

	if current != nil {
		textBlocks = append(textBlocks, strings.Join(current, "\n"))
		current = nil
	}

	return textBlocks
}

func BlocksToHTML(textBlocks []string) string {
	var nodes []HTMLNode
	for _, textBlock := range textBlocks {
		var node HTMLNode

		switch bt := getBlockType(textBlock); bt {
		case PARAGRAPH:
			node = createParagraph(textBlock)
		case CODE:
			node = createCode(textBlock)
		case HEADING:
			node = createHeading(textBlock)
		case QUOTE:
			node = createQuote(textBlock)
		case UNORDERED_LIST:
			node = createUnorderedList(textBlock)
		case ORDERED_LIST:
			node = createOrderedList(textBlock)
		default:
			log.Panic(fmt.Sprint("converting md to block - unexpected block type: ", bt))
		}

		nodes = append(nodes, node)
	}

	root := NewParentNode("div", "", nodes)
	return root.ToHTML()
}

func createHeading(text string) HTMLNode {
	idx := strings.Index(text, " ")
	tag := fmt.Sprintf("h%d", idx)
	value := text[idx+1:]

	return NewLeafNode(tag, "", value)
}

func createParagraph(text string) HTMLNode {
	return NewLeafNode("p", "", text)
}

func createCode(text string) HTMLNode {
	code := NewLeafNode("code", "", text[3:len(text)-3])
	return NewParentNode("pre", "", []HTMLNode{code})
}

func createQuote(text string) HTMLNode {
	var value []string
	for _, line := range strings.Split(text, "\n") {
		value = append(value, line[2:])
	}
	return NewLeafNode("blockquote", "", strings.Join(value, "\n"))
}

func createUnorderedList(text string) HTMLNode {
	var listItems []HTMLNode
	for _, line := range strings.Split(text, "\n") {
		item := NewLeafNode("li", "", line[2:])
		listItems = append(listItems, item)
	}
	return NewParentNode("ul", "", listItems)
}

func createOrderedList(text string) HTMLNode {
	var listItems []HTMLNode
	for _, line := range strings.Split(text, "\n") {
		item := NewLeafNode("li", "", line[3:])
		listItems = append(listItems, item)
	}
	return NewParentNode("ol", "", listItems)
}
