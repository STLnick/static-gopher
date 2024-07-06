package parse

import "fmt"

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

func NewTextNode(text string, textType TextType, url string) TextNode {
	return TextNode{
		text:     text,
		textType: textType,
		url:      url,
	}
}
