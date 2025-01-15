package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ktr0731/go-fuzzyfinder"
)

//go:embed amharic_bible.json
var f embed.FS

type Track struct {
	Name      string `json:"name,omitempty"`
	AlbumName string `json:"album_name,omitempty"`
	Artist    string `json:"artist,omitempty"`
}

type Chapter struct {
	Chapter string   `json:"chapter,omitempty"`
	Title   string   `json:"title,omitempty"`
	Verses  []string `json:"verses,omitempty"`
}
type Book struct {
	Title    string    `json:"title,omitempty"`
	Abrv     string    `json:"abrv,omitempty"`
	Chapters []Chapter `json:"chapters,omitempty"`
}

type Bible struct {
	Books []Book `json:"books"`
}

type Verse struct {
	BookTitle string `json:"book_title,omitempty"`
	verse     string
	Chapter   string `json:"chapter,omitempty"`
	Number    int    `json:"number,omitempty"`
}

func buildSearchableDocument(bible Bible) []Verse {
	var content []Verse
	for _, book := range bible.Books {
		for _, chapter := range book.Chapters {
			for num, verse := range chapter.Verses {
				entry := Verse{BookTitle: book.Title, verse: verse, Chapter: chapter.Chapter, Number: num + 1}
				content = append(content, entry)
			}
		}
	}
	return content
}

func main() {
	file, err := f.ReadFile("amharic_bible.json")
	if err != nil {
		log.Fatal(err)
	}
	var bible Bible
	if err := json.Unmarshal(file, &bible); err != nil {
		log.Fatal(err)
	}
	content := buildSearchableDocument(bible)
	idx, err := fuzzyfinder.FindMulti(
		content,
		func(i int) string {
			return content[i].verse
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Verse: %s %s:%d", content[i].BookTitle, content[i].Chapter, content[i].Number)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	selectedVerse := content[idx[0]]
	p := fmt.Sprintf("%v %s %s: %d", selectedVerse.verse, selectedVerse.BookTitle, selectedVerse.Chapter, selectedVerse.Number)
	fmt.Println(p)
}
