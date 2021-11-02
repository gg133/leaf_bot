.PHONY:
.SILENT:

build:
	go build -o ./.bin/leaf_bot main.go

run: build
	./.bin/leaf_bot

build-image:
	docker build -t leaf_bot:v0.1 .

start-container:
	docker run --name leaf_bot -p 80:80 --env-file .env leaf_bot:v0.1