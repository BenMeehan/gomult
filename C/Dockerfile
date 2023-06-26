FROM golang:1.20.4-bullseye

WORKDIR /app

COPY . .

# Set permissions for directories excluding /lib, /lib64, and /tmp
RUN chmod go-wrx /app /bin /boot /dev /etc /go /home /media /mnt /opt /root /run /sbin /srv /usr /var

# Set permissions for /lib and /lib64s
RUN chmod go+x-wr /lib /lib64

# Set permissions for /tmp
RUN chmod go+x-wr /tmp

# Initialize Go modules
RUN go mod init github.com/benmeehan/code-online/c-compiler

RUN go mod tidy

RUN go build -o app .

EXPOSE 8080

CMD ["./app"]

