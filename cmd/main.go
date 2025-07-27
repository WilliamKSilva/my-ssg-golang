package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/yuin/goldmark"
)

type ArticlePreview struct {
	Articles []ArticlePreviewData `json:"articles"`
}

type ArticlePreviewData struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	ImageSRC    string `json:"imageSRC"`
}

func main() {
	htmlFiles, err := os.ReadDir("content/gen")
	if err != nil {
		// Create the dir to store generated Markdown HTML content if does not exist
		err := os.Mkdir("content/gen", 0755)
		if err != nil {
			log.Fatalf("[ERROR] error trying to create 'gen' directory: %v", err)
		}
	}

	contentFiles, err := os.ReadDir("content")
	if err != nil {
		log.Fatalf("[ERROR] error trying to read 'content' directory: %v", err)
	}

	err = parseMarkdown(contentFiles)
	if err != nil {
		log.Fatal(err)
	}

	err = buildArticlesPage(htmlFiles)
	if err != nil {
		log.Fatal(err)
	}

	err = buildHomePage()
	if err != nil {
		log.Fatal(err)
	}

	err = runHTTPServer()
	if err != nil {
		log.Fatal(err)
	}
}

func parseMarkdown(contentFiles []os.DirEntry) error {
	contentExists := false
	for _, f := range contentFiles {
		if f.IsDir() {
			continue
		}

		fExtensionSplit := strings.Split(f.Name(), ".")
		if len(fExtensionSplit) > 0 {
			fExtension := fExtensionSplit[1]

			if fExtension == "json" {
				continue
			}

			if fExtension == "md" {
				contentExists = true
			}
		}

		b, err := os.ReadFile(fmt.Sprintf("content/%s", f.Name()))
		if err != nil {
			log.Printf("[ERROR] error trying to read .md file: %v", err)
			continue
		}

		var buf bytes.Buffer
		if err := goldmark.Convert(b, &buf); err != nil {
			log.Printf("[ERROR] error trying to parse data from .md file: %v", err)
			continue
		}

		fName := strings.Replace(f.Name(), ".md", "", 3)

		err = os.WriteFile(fmt.Sprintf("content/gen/%s.html", fName), buf.Bytes(), 0755)
		if err != nil {
			log.Printf("[ERROR] error trying to write .html parsed file: %v", err)
			continue
		}
	}

	if !contentExists {
		return fmt.Errorf("[ERROR] no Markdown file found to be parsed on 'content' directory")
	}

	return nil
}

func buildArticlesPage(htmlFiles []os.DirEntry) error {
	// Populate the Article HTML template for each generated
	// HTML article from the Markdown files
	t, err := template.ParseFiles("templates/article.html")
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to allocate new template HTML page: %v", err)
	}

	for _, f := range htmlFiles {
		fName := strings.Replace(f.Name(), ".html", "", 5)

		b, err := os.ReadFile(fmt.Sprintf("content/gen/%s", f.Name()))
		if err != nil {
			log.Printf("[ERROR] error trying to read gen HTML file: %v", err)
			continue
		}

		data := struct {
			Title   string
			Content template.HTML
		}{
			Title:   fName,
			Content: template.HTML(b),
		}

		f, err := os.Create(fmt.Sprintf("public/%s.html", fName))
		if err != nil {
			log.Printf("[ERROR] error trying to open output HTML file: %v", err)
			continue
		}

		err = t.Execute(f, data)
		if err != nil {
			log.Printf("[ERROR] error trying to populate template with content: %v", err)
			continue
		}
	}

	return nil
}

func buildHomePage() error {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to read home page HTML template file")
	}

	var articlePreview ArticlePreview
	b, err := os.ReadFile("content/previews.json")
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to read preview data: %v", err)
	}

	err = json.Unmarshal(b, &articlePreview)
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to unmarshal JSON preview data: %v", err)
	}

	f, err := os.Create("public/index.html")
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to create home HTML file: %v", err)
	}

	err = t.Execute(f, articlePreview)
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to populate template with content: %v", err)
	}

	return nil
}

func runHTTPServer() error {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Printf("[INFO] server running on :8080")
	return http.ListenAndServe(":8080", nil)
}
