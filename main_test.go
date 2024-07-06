package main

import (
    "testing"
    "strings"
)

type TextToTextNodesTest struct {
    text string
    expectedLen int
    typeCounts map[TextType]int
}

var textToTextNodeTests = []TextToTextNodesTest{
    TextToTextNodesTest{
        "This has *italic text*, **bold** text, and `a little bit of code`.",
        7,
        map[TextType]int{
            TEXT: 4,
            BOLD: 1,
            ITALIC: 1,
            CODE_TEXT: 1,
        },
    },
    TextToTextNodesTest{
        "Nothing special",
        1,
        map[TextType]int{
            TEXT: 1,
            BOLD: 0,
            ITALIC: 0,
            CODE_TEXT: 0,
        },
    },
    TextToTextNodesTest{
        "I *have* **many** *blocks* of **stylized text** `within` this `one string` with **bolded**",
        14,
        map[TextType]int{
            TEXT: 7,
            BOLD: 3,
            ITALIC: 2,
            CODE_TEXT: 2,
        },
    },
    TextToTextNodesTest{
        "[i am link](https://www.google.com)",
        1,
        map[TextType]int{
            LINK: 1,
        },
    },
    TextToTextNodesTest{
        "![i am image](https://picsum.photos/200/300)",
        1,
        map[TextType]int{
            IMAGE: 1,
        },
    },
}

func TestTextToTextNodes(t *testing.T) {
    for _, test := range textToTextNodeTests {
        textNodes := TextToTextNodes(test.text)
        if len(textNodes) != test.expectedLen {
            t.Errorf("Expected (%d) nodes but got (%d)\n", test.expectedLen, len(textNodes))
        }
        counts := map[TextType]int{
            TEXT: 0,
            BOLD: 0,
            ITALIC: 0,
            CODE_TEXT: 0,
            LINK: 0,
            IMAGE: 0,
        }
        for _, node := range textNodes {
            counts[node.textType] += 1
        }
        for key, value := range counts {
            if value != test.typeCounts[key] {
                t.Errorf("Expected (%d) '%s' nodes but got (%d)\n", test.typeCounts[key], key.Name(), value)
            }
        }
    }
}

type MarkdownToTextBlockTest struct {
    md string
    expectedBlocks int
}

var markdownToTextBlockTests = []MarkdownToTextBlockTest{
    MarkdownToTextBlockTest{
        "# heading1\n\ntext\n\n## heading2\n\n\n> q1\n> q2\n\n- li1\n- li2\n\n1. li1\n2. li2",
        6,
    },
    MarkdownToTextBlockTest{
        "```\nconst x = 2;\n```\n\n## heading 2\n\njust some text",
        3,
    },
    MarkdownToTextBlockTest{
        "# heading1\n\nthis text has an ![image](https://picsum.photos/300/200) and a [link](https://www.google.com)",
        2,
    },
}

func TestMarkdownToTextBlocks(t *testing.T) {
    for _, test := range markdownToTextBlockTests {
        results := MarkdownToBlocks(test.md)
        if len(results) != test.expectedBlocks {
            t.Errorf("Expected (%d) blocks but got (%d)", test.expectedBlocks, len(results))
        }
    }
}

type BlocksToHTMLTest struct {
    md string
    expectedTags []string
}

var blocksToHTMLTests = []BlocksToHTMLTest{
    BlocksToHTMLTest{
        "# heading1\n\ntext\n\n## heading2\n\n\n> q1\n> q2\n\n- li1\n- li2\n\n1. li1\n2. li2",
        []string{"div", "h1", "p", "h2", "blockquote", "ul", "ol", "li"},
    },
    BlocksToHTMLTest{
        "```\nconst x = 2;\n```\n\n## heading 2\n\njust some text",
        []string{"div", "pre", "code", "h2", "p"},
    },
    BlocksToHTMLTest{
        "# heading1\n\nthis text has an ![image](https://picsum.photos/300/200) and a [link](https://www.google.com)",
        []string{"div", "h1", "img", "a", "p"},
    },
}

func TestBlocksToHTML(t *testing.T) {
    for _, test := range blocksToHTMLTests {
        blocks := MarkdownToBlocks(test.md)
        result := BlocksToHTML(blocks)
        for _, tag := range test.expectedTags {
            if strings.Index(result, tag) == -1 {
                t.Errorf("Expected (%s) tag but was not present", tag)
            }
        }
    }
}


