CONTENT = content
PUBLIC_PATH = public

clean-html:
	@echo "Deleting generated HTML..."
	@rm $(CONTENT)/gen/*.html
	@rm $(PUBLIC_PATH)/*.html
	@echo "Sucessfully deleted HTML files..."