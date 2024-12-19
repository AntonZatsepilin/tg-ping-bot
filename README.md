# tg-ping-bot

### Launch

Create an .env file like this in the root directory of your project:

``` .env
TELEGRAM_BOT_TOKEN=...
TELEGRAM_CHAT_ID=...
```

Change these parameters in the docker-compose.yml to suit yourself:

``` yml
INTERVAL: ...
REQUEST_TIMEOUT: ...
WORKERS_COUNT: ...
```

Add links in the links.yml to suit yourself:

``` yml
links:
  - ...
  - ...
  - ...
  ...
```

Enter this command in a terminal running Docker:

```
docker-compose up --build
```
