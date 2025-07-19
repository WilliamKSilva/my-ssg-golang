package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/yuin/goldmark"
)

const tpl = `
	<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="UTF-8" />
				<title>{{ .Title }}</title>
				<link rel="stylesheet" href="./assets/style.css" />
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

func parseMarkdown() error {
	dirs, err := os.ReadDir("content")
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to read markdown files from dir: %v", err)
	}

	for _, d := range dirs {
		if d.IsDir() {
			continue
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

func main() {
	err := parseMarkdown()
	if err != nil {
		log.Fatal(err)
	}

	dirs, err := os.ReadDir("content/gen")
	if err != nil {
		log.Fatalf("[ERROR] error trying to read gen HTML files: %v", err)
	}

	for _, d := range dirs {
		t, err := template.New("webpage").Parse(tpl)
		if err != nil {
			log.Fatalf("[ERROR] error trying to allocate new template HTML page: %v", err)
		}

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
}
