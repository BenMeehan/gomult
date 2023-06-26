FROM golang:1.20.4-bullseye

WORKDIR /app

COPY . .

# Set permissions for directories excluding /lib, /lib64, and /tmp
RUN chmod go-wrx /bin /boot /dev /etc /go /home /media /mnt /opt /run /sbin /srv /usr /var

# Set permissions for /lib and /lib64
RUN chmod go+x-wr /lib /lib64 /usr

# Set permissions for /root
RUN chmod go+x-rw /root

# Set permissions for /app
RUN chmod go+x-rw /app

# Go expects bins in /usr/local. Remove all other permissions
RUN chmod go-x /usr/bin /usr/lib /usr/games /usr/include /usr/libexec /usr/sbin /usr/share /usr/src

# Initialize Go modules
RUN go mod init github.com/benmeehan/code-online/go-compiler

RUN go mod tidy

RUN go build -o app .

EXPOSE 8080

CMD ["./app"]