FROM golang:latest

WORKDIR /app

# Install required libraries for C++ compilation
RUN apt-get update && \
    apt-get install -y build-essential && \
    apt-get install -y libstdc++6

COPY . .

# Create a restricted user
RUN adduser restricteduser

# Set permissions for directories excluding /lib, /lib64, and /tmp
RUN chmod go-wrx /app /bin /boot /dev /etc /go /home /media /mnt /opt /proc /root /run /sbin /srv /usr /var

# Set permissions for /lib and /lib64
RUN chmod go+x-wr /lib /lib64

# Set permissions for /tmp
RUN chmod go+x-wr /tmp

# Initialize Go modules
RUN go mod init github.com/benmeehan/code-online/c-compiler

RUN go mod tidy

RUN go build -o app .

EXPOSE 8080

CMD ["./app"]

