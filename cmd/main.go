package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
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

const tpl = `
	<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="UTF-8" />
				<title>{{ .Title }}</title>
				<link rel="stylesheet" href="./assets/articles.css" />
			</head>
		<body>
			<header>
				<a href="/">← Home</a>
			</header>
			<main>
				{{ .Content }}
			</main>
			<footer>
				<p>© 2025 Your Name</p>
			</footer>
		</body>
	</html>
`

func main() {
	err := parseMarkdown()
	if err != nil {
		log.Fatal(err)
	}

	err = buildArticlesPage()
	if err != nil {
		log.Fatal(err)
	}

	err = buildHomePage()
	if err != nil {
		log.Fatal(err)
	}
}

func parseMarkdown() error {
	dirs, err := os.ReadDir("content")
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to read markdown files from dir: %v", err)
	}

	for _, d := range dirs {
		if d.IsDir() {
			continue
		}

		fExtensionSplit := strings.Split(d.Name(), ".")
		if len(fExtensionSplit) > 0 {
			fExtension := fExtensionSplit[1]

			if fExtension == "json" {
				continue
			}
		}

		b, err := os.ReadFile(fmt.Sprintf("content/%s", d.Name()))
		if err != nil {
			log.Printf("[ERROR] error trying to read .md file: %v", err)
			continue
		}

		var buf bytes.Buffer
		if err := goldmark.Convert(b, &buf); err != nil {
			log.Printf("[ERROR] error trying to parse data from .md file: %v", err)
			continue
		}

		fName := strings.Replace(d.Name(), ".md", "", 3)

		err = os.WriteFile(fmt.Sprintf("content/gen/%s.html", fName), buf.Bytes(), 0755)
		if err != nil {
			log.Printf("[ERROR] error trying to write .html parsed file: %v", err)
			continue
		}
	}

	return nil
}

func buildArticlesPage() error {
	dirs, err := os.ReadDir("content/gen")
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to read gen HTML files: %v", err)
	}

	// Populate the Article HTML template for each generated
	// HTML article from the Markdown files
	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to allocate new template HTML page: %v", err)
	}

	for _, d := range dirs {
		fName := strings.Replace(d.Name(), ".html", "", 5)

		b, err := os.ReadFile(fmt.Sprintf("content/gen/%s", d.Name()))
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
