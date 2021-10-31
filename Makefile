.PHONY:
.SILENT:

build:
	go build -o ./.bin/leaf_bot main.go

run: build
	./.bin/leaf_bot
