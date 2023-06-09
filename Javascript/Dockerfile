FROM golang:latest

WORKDIR /app

# Install Node.js
RUN apt-get update && \
    apt-get install -y nodejs && \
    apt-get clean

COPY . .

# Set permissions for directories excluding /lib, /lib64, and /tmp
RUN chmod go-wrx /app /bin /boot /dev /etc /go /home /media /mnt /opt /root /run /sbin /srv /usr /var

# Set permissions for /usr
RUN chmod go+x-wr /usr

# Set permissions for /tmp
RUN chmod go+x-wr /tmp

# Node.js bin is in /usr/bin and shared libs in /usr/lib. Remove all other permissions
RUN chmod go-x /usr/games /usr/include /usr/libexec /usr/local /usr/sbin /usr/share /usr/src

# Initialize Go modules
RUN go mod init github.com/benmeehan/code-online/javascript-compiler

RUN go mod tidy

RUN go build -o app .

EXPOSE 8080

CMD ["./app"]