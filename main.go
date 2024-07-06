package main

import (
    "github.com/STLnick/static-gopher/parse"
	"fmt"
	"log"
	"os"
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

    testFilePath := "test.md"
    md, err := os.ReadFile(testFilePath)
    if err != nil {
        log.Panic(fmt.Sprint("reading markdown file:", err))
    }
    fmt.Println("read markdown into memory.")

    blocks := parse.MarkdownToBlocks(string(md))
    fmt.Println("created blocks from markdown.")
    
    html := parse.BlocksToHTML(blocks)
    fmt.Println("created HTML from blocks.")
    
    err = os.WriteFile("test.html", []byte(html), 0644)
    if err != nil {
        log.Panic(fmt.Sprint("writing html file:", err))
    }
    fmt.Println("HTML written to file")
}

