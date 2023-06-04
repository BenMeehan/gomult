FROM golang:latest

WORKDIR /app

# Install Java Development Kit (JDK)
RUN apt-get update && \
    apt-get install -y default-jdk

COPY . .

# Set permissions for directories excluding /lib, /lib64, and /tmp
RUN chmod go-wrx /app /bin /boot /dev /etc /go /home /media /mnt /opt /proc /root /run /sbin /srv /usr /var

# Set permissions for /lib and /lib64
RUN chmod go+x-wr /lib /lib64 /usr /etc

# Set permissions for /tmp
RUN chmod go+x-wr /tmp

# Java bin is in /usr/bin and expects shared libs in /usr/lib. Remove all other permissions
RUN chmod go-x /usr/games /usr/include /usr/libexec /usr/local /usr/sbin /usr/share /usr/src

# Initialize Go modules
RUN go mod init github.com/benmeehan/code-online/java-compiler

RUN go mod tidy

RUN go build -o app .

EXPOSE 8080

CMD ["./app"]