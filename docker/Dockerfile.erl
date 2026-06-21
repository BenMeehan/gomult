FROM golang:1.26-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o gomult .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends erlang \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/gomult /usr/local/bin/gomult
COPY configs/erl.yaml /etc/gomult/config.yaml
RUN useradd -m -u 1000 runner
USER 1000
EXPOSE 8080
CMD ["gomult", "--config", "/etc/gomult/config.yaml"]
