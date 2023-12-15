.PHONY: build db sql

build:
	@go build -o gen

sql: build
	@./gen sql

db: build
	@./gen db