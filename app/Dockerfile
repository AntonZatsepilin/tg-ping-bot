FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY ./schema /schema

RUN apk add --no-cache curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin

RUN go build -o main ./cmd

ENV TELEGRAM_BOT_TOKEN="your_default_bot_token"
ENV TELEGRAM_CHAT_ID="your_default_chat_id"

CMD ["./main"]
