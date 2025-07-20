# Introduction
- This repo stores a implementation of a basic SSG (Static Site Generator) using Golang to build a personal blog. 

- The idea is to generate HTML content for articles wrote in `Markdown` files.

- This project was built for personal use, it has a "specific" usecase, but maybe if you also need a very basic
personal blog, where the goal is to store blog posts you want to write about you can use this implementation.

- This approach makes easier to deploy static files to services like Github Pages or Netlify, for example.

# How to use

### Blog Posts
- When I'am studying a subject, I like to take notes and write about this particular thing on a `Markdown` file, so this is the file supported in this implementation.
- Every .md file you place inside the `content` directory will be processed when you run `make run` and be turned into a new HTML file that you can access through the URL: public/file.html.
- To make the list of blog posts on the home page to work you need to populate the file `content/previews.json`.

### Templates
- The `templates` directory stores the core structure of the two main webpages that our blog will have: the home page `index.html` and the articles webpage `article.html`.
- If you want to change the core of the two main webpages you will need to make changes to the HTML template files. Example: change your name on the footer of the webpages.

### Styling
- You can change the layout and styling of the generated webpages through the CSS files inside `public/assets`.

### Run
- You need `Golang` toolchain installed: https://go.dev/doc/install.
- Place your blog posts inside the `content` directory.
- Run `make run`.
- Now you have your static pages inside the `public` directory.