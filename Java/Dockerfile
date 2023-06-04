FROM golang:latest

WORKDIR /app

# Install Java Development Kit (JDK)
RUN apt-get update && \
    apt-get install -y default-jdk

COPY . .

# Create a restricted user
RUN adduser restricteduser

# Set permissions for directories excluding /lib, /lib64, and /tmp
RUN chmod go-wrx /app /bin /boot /dev /etc /go /home /media /mnt /opt /proc /root /run /sbin /srv /usr /var

# Set permissions for /lib and /lib64
RUN chmod go+x-wr /lib /lib64 /usr

# Set permissions for /tmp
RUN chmod go+x-wr /tmp

# Initialize Go modules
RUN go mod init github.com/benmeehan/code-online/c-compiler

RUN go mod tidy

RUN go build -o app .

EXPOSE 8080

USER restricteduser

CMD ["./app"]