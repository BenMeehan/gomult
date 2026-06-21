FROM golang:1.26-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o gomult .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends curl unzip default-jre \
    && curl -sL "https://github.com/JetBrains/kotlin/releases/download/v1.9.22/kotlin-compiler-1.9.22.zip" -o /tmp/kotlin.zip \
    && unzip -q /tmp/kotlin.zip -d /usr/local/kotlin \
    && ln -s /usr/local/kotlin/kotlinc/bin/kotlin /usr/local/bin/kotlin \
    && ln -s /usr/local/kotlin/kotlinc/bin/kotlinc /usr/local/bin/kotlinc \
    && rm /tmp/kotlin.zip \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/gomult /usr/local/bin/gomult
COPY configs/kt.yaml /etc/gomult/config.yaml
RUN useradd -m -u 1000 runner
USER 1000
EXPOSE 8080
CMD ["gomult", "--config", "/etc/gomult/config.yaml"]
