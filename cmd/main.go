package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/yuin/goldmark"
)

type PostPreview struct {
	Posts []PostPreviewData `json:"articles"`
}

type PostPreviewData struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	ImageSRC    string `json:"imageSRC"`
}

var httpServer bool
var outputPath string

func main() {
	flag.BoolVar(&httpServer, "httpServer", false, "Start a HTTP server with the generated files served so the content can be visualized")
	flag.StringVar(&outputPath, "outputPath", "public/", "Specify the output path of the final HTML content")

	flag.Parse()

	// Create the dir to store generated Markdown HTML content if does not exist
	_, err := os.ReadDir("content/gen")
	if err != nil {
		err := os.Mkdir("content/gen", 0755)
		if err != nil {
			log.Fatalf("[ERROR] error trying to create 'gen' directory: %v", err)
		}
	}

	// Create the dir to store output HTML static files
	_, err = os.ReadDir(outputPath)
	if err != nil {
		err := os.Mkdir(outputPath, 0755)
		if err != nil {
			log.Fatalf("[ERROR] error trying to create '%s' directory: %v", outputPath, err)
		}
	}

	// Create the dir to store output posts HTML static files
	blogPostsPath := fmt.Sprintf("%s/posts", outputPath)
	_, err = os.ReadDir(blogPostsPath)
	if err != nil {
		err := os.Mkdir(blogPostsPath, 0755)
		if err != nil {
			log.Fatalf("[ERROR] error trying to create '%s' directory: %v", blogPostsPath, err)
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

	err = buildPostPages()
	if err != nil {
		log.Fatal(err)
	}

	err = buildIndexPage()
	if err != nil {
		log.Fatal(err)
	}

	err = os.CopyFS(fmt.Sprintf("%s/assets/", outputPath), os.DirFS("assets"))
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		log.Fatalf("[ERROR] error tring to copy assets to final output path: %v", err)
	}

	if httpServer {
		err = runHTTPServer()
		if err != nil {
			log.Fatal(err)
		}
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

func buildPostPages() error {
	// Populate the Post HTML template for each generated
	// HTML post from the Markdown files
	t, err := template.ParseFiles("templates/post.html")
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to allocate new template HTML page: %v", err)
	}

	htmlFiles, err := os.ReadDir("content/gen")
	if err != nil {
		err := os.Mkdir("content/gen", 0755)
		if err != nil {
			log.Fatalf("[ERROR] error trying to create 'gen' directory: %v", err)
		}
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

		f, err := os.Create(fmt.Sprintf("%s/posts/%s.html", outputPath, fName))
		if err != nil {
			log.Printf("[ERROR] error trying to create output HTML file: %v", err)
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

func buildIndexPage() error {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to read home page HTML template file")
	}

	var postPreview PostPreview
	b, err := os.ReadFile("content/previews.json")
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to read preview data: %v", err)
	}

	err = json.Unmarshal(b, &postPreview)
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to unmarshal JSON preview data: %v", err)
	}

	f, err := os.Create(fmt.Sprintf("%s/index.html", outputPath))
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to create home HTML file: %v", err)
	}

	err = t.Execute(f, postPreview)
	if err != nil {
		return fmt.Errorf("[ERROR] error trying to populate template with content: %v", err)
	}

	return nil
}

func runHTTPServer() error {
	fs := http.FileServer(http.Dir(outputPath))
	http.Handle("/", fs)

	log.Printf("[INFO] server running on :8080")
	return http.ListenAndServe(":8080", nil)
}
