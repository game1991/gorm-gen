.PHONY: build db sql

build:
	@rm -rf gen
	@go build -o gen

sql: build
	@./gen sql

db: build
	@./gen db