// Command read prints the metadata and content of an EPUB file passed as the
// first command-line argument.
//
// Usage:
//
//	go run ./examples/read book.epub
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/raitucarp/epub"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: read <file.epub>")
	}

	book, err := epub.OpenReader(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Title:", strings.Join(book.Title(), ", "))
	fmt.Println("Author:", strings.Join(book.Author(), ", "))
	fmt.Println("Language:", strings.Join(book.Language(), ", "))
	fmt.Println("Identifier:", book.UID())
	fmt.Println("Version:", book.Version())

	for _, id := range book.ListContentDocumentIds() {
		md := book.ReadContentMarkdownById(id)
		fmt.Printf("\n--- %s ---\n%s\n", id, md)
	}
}
