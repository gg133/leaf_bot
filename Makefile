.PHONY:
.SILENT:

build:
	go build -o ./.bin/leaf_bot main.go

run: build
	sudo setcap 'cap_net_bind_service=+ep' ~/Projects/leaf_bot/.bin/leaf_bot
	./.bin/leaf_bot
	
