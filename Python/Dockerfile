FROM golang:latest

WORKDIR /app

# Install Python and required libraries
RUN apt-get update && \
    apt-get install -y python3 && \
    apt-get clean

COPY . .

# Create a restricted user
RUN adduser restricteduser

# Set permissions for directories excluding /lib, /lib64, and /tmp
RUN chmod go-wrx /app /bin /boot /dev /etc /go /home /media /mnt /opt /proc /root /run /sbin /srv /usr /var

# Set permissions for /lib and /lib64
RUN chmod go+x-wr /lib /lib64 

# Set permissions for /tmp
RUN chmod go+x-wr /tmp

# Python expects shared libs in /usr/lib. Remove all other permissions
# RUN chmod go-x /usr/bin /usr/games /usr/include /usr/libexec /usr/local /usr/sbin /usr/share /usr/src

# Initialize Go modules
RUN go mod init github.com/benmeehan/code-online/python-compiler

RUN go mod tidy

RUN go build -o app .

EXPOSE 8080

CMD ["./app"]