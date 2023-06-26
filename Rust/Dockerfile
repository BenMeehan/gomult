FROM golang:1.20.4-bullseye

WORKDIR /app

# Install required libraries for Rust compilation
RUN apt-get update && \
    apt-get install -y build-essential && \
    apt-get install -y rustc

COPY . .

# Set permissions for directories excluding /lib, /lib64, and /tmp
RUN chmod go-wrx /app /bin /boot /dev /etc /go /home /media /mnt /opt /root /run /sbin /srv /usr /var

# Set permissions for /lib and /lib64
RUN chmod go+x-wr /usr

# Set permissions for /tmp
RUN chmod go+x-wr /tmp

# Rust expects shared libs in /usr/lib. Remove all other permissions
RUN chmod go-x /usr/bin /usr/games /usr/include /usr/libexec /usr/local /usr/sbin /usr/share /usr/src

# Initialize Go modules
RUN go mod init github.com/benmeehan/code-online/rust-compiler

RUN go mod tidy

RUN go build -o app .

EXPOSE 8080

CMD ["./app"]