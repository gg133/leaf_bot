.PHONY:
.SILENT:

build:
	go build -o ./.bin/leaf_bot main.go

run: build
	sudo setcap 'cap_net_bind_service=+ep' ~/Projects/leaf_bot/.bin/leaf_bot
	./.bin/leaf_bot

build-image:
	sudo docker build -t leaf_bot:v0.1 .

start-container:
	sudo docker run --name leaf_bot -p 80:80 --env-file .env leaf_bot:v0.1