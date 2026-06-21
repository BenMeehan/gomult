FROM golang:1.26-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o gomult .

FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential pkg-config protobuf-compiler \
    libprotobuf-dev libnl-3-dev libnl-route-3-dev git ca-certificates \
    flex curl unzip \
    gcc g++ gfortran gdc golang rustc \
    default-jdk scala groovy \
    python3 nodejs npm ruby php-cli \
    perl lua5.4 \
    ghc ocaml \
    elixir erlang \
    r-base racket swi-prolog \
    mono-devel nasm fp-compiler \
    libnl-3-200 libnl-route-3-200 \
    && npm install -g typescript ts-node \
    && rm -rf /root/.npm /var/lib/apt/lists/*

RUN curl -sL "https://github.com/JetBrains/kotlin/releases/download/v1.9.22/kotlin-compiler-1.9.22.zip" -o /tmp/kotlin.zip \
    && unzip -q /tmp/kotlin.zip -d /usr/local/kotlin \
    && ln -s /usr/local/kotlin/kotlinc/bin/kotlin /usr/local/bin/kotlin \
    && ln -s /usr/local/kotlin/kotlinc/bin/kotlinc /usr/local/bin/kotlinc \
    && rm /tmp/kotlin.zip

RUN git clone --depth 1 https://github.com/google/nsjail.git /tmp/nsjail \
    && cd /tmp/nsjail && make -j"$(nproc)" \
    && cp nsjail /usr/local/bin/ \
    && rm -rf /tmp/nsjail \
    && apt-get purge -y pkg-config protobuf-compiler libprotobuf-dev git curl unzip flex \
    && apt-get autopurge -y \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/gomult /usr/local/bin/gomult
COPY config.yaml /etc/gomult/config.yaml

RUN useradd -m -u 1000 -s /bin/bash runner

USER 1000
EXPOSE 8080

CMD ["gomult", "--config", "/etc/gomult/config.yaml"]
