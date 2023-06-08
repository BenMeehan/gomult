FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod init github.com/benmeehan111/lb

RUN go mod tidy

RUN go build -o app

EXPOSE 8080

CMD ["./app"]
