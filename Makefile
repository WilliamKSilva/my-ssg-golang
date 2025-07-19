CONTENT = content
PUBLIC_PATH = public

run:
	@echo "Running SSG with Go toolchain"
	go run cmd/main.go

clean-html:
	@echo "Deleting generated HTML..."
	@rm $(CONTENT)/gen/*.html
	@rm $(PUBLIC_PATH)/*.html
	@echo "Sucessfully deleted HTML files..."