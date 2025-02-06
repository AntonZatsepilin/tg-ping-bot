FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd

ENV TELEGRAM_BOT_TOKEN="your_default_bot_token"
ENV TELEGRAM_CHAT_ID="your_default_chat_id"

CMD ["./main"]
