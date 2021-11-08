# Leaf (pocket) bot

LeafBot - is a Telegram bot, which allows you to save link in the application [Pocket](https://getpocket.com/). Leafbot is a client for Pocket. 

## Used libraries:
- The [telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) is used to work with the telegram bot. 
- [BoltDB](https://github.com/boltdb/bolt) is used as storage.
- [Viper](https://github.com/spf13/viper) is used to get data from config files

## Implement
First at all, when the bot started, bot generates the link with unique request token, which created and stored in DB when user started the bot, and send it to the user, at the click of the link authorization process starts.

To implement user authorization, together with the bot, an HTTP server is started on port 80. HTTP server is using for creating access token, put it in the DB and redirecting from Pocket after successful user authorization.

After the authorization process, user can save links to Pocket by LeafBot.

## How to use
For local using you can clone this repository via git:

```bash
git clone github.com/yalagtyarzh/leaf_bot
```

or install it via go: 

```bash
go install github.com/yalagtyarzh/leaf_bot
```

After, you need to:
- Build a program with the pre-registred enviroment dependencies using [Docker](https://www.docker.com/):

```bash
docker build -t <build-name> .
```

and run it like that:

```bash
docker run -e TOKEN -e CONSUMER_KEY -e AUTH_SERVER_URL -d \ --name <container-name>
```

You can use Makefile for building docker image and running docker container (Makefile must be installed).

- Uncomment os.Setenv functions in [config.go](https://github.com/yalagtyarzh/leaf_bot/blob/master/config/config.go) and enter telegram bot token instead of YourAwesomeBotToken, enter Pocket consumer key instead of YourAwesomePocketToken and enter vaild authorization server instead and build program using go tools or Makefile (Makefile must be installed) (not recommended).

```bash
make run
```

After these steps, the bot should start (hopefully Kappa).

## Stack:
- [Go](https://golang.org/)
- [Docker](https://www.docker.com/)
