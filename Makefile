## docs/new page=PAGE_NAME.md: Creates a new docs page in the docs/content/docs directory
.PHONY: docs/new
docs/new:
ifndef page
	@echo "Missing page name. Use docs/new page=PAGE_NAME.md"
	@exit 1
endif
	hugo new docs/$(page) --source docs

## docs/serve: Builds and hosts the docs locally
.PHONY: docs/serve
docs/serve:
	hugo serve --buildDrafts --source docs --openBrowser

# Utilites
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
