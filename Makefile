OUTPUT_PATH=blog

run:
	@echo "Running SSG with Go toolchain"
	go run cmd/main.go -outputPath=$(OUTPUT_PATH)

run-http-server:
	@echo "Running SSG with Go toolchain"
	go run cmd/main.go -outputPath=$(OUTPUT_PATH) -httpServer=true

clean:
	@echo "Deleting generated HTML..."
	@rm -rf content/gen
	@rm -rf $(OUTPUT_PATH)
	@echo "Sucessfully deleted HTML files..."