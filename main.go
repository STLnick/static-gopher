package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
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
	for _, tn := range textNodes {
		fmt.Println(tn)
	}
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
	nodes, err = splitNodesLinks(nodes)
	if err != nil {
		log.Panic("splitting on links:", err)
	}
	nodes, err = splitNodesImages(nodes)
	if err != nil {
		log.Panic("splitting on images:", err)
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
	// TODO
	return nodes, nil
}

func splitNodesImages(nodes []TextNode) ([]TextNode, error) {
	// TODO
	return nodes, nil
}
