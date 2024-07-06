package main

import (
    "testing"
)

type TextToTextNodesTest struct {
    text string
    expectedLen int
    typeCounts map[TextType]int
}

var tests = []TextToTextNodesTest{
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
}

func TestTextToTextNodes(t *testing.T) {
    for _, test := range tests {
        textNodes := TextToTextNodes(test.text)
        if len(textNodes) != test.expectedLen {
            for i, n := range textNodes {
                t.Log(i, n)
            }
            t.Errorf("Expected (%d) nodes but got (%d)\n", test.expectedLen, len(textNodes))
        }
        counts := map[TextType]int{

            TEXT: 0,
            BOLD: 0,
            ITALIC: 0,
            CODE_TEXT: 0,
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
