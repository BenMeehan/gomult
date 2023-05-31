FROM golang:latest

WORKDIR /app

# Install Python and required libraries
RUN apt-get update && \
    apt-get install -y python3 && \
    apt-get clean

COPY . .

# Initialize Go modules
RUN go mod init github.com/benmeehan/code-online/python-compiler

RUN go mod tidy

RUN go build -o app .

EXPOSE 8080

CMD ["./app"]