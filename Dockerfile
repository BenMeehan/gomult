FROM golang:1.19.3-bullseye

WORKDIR "/app"

RUN apt update

# Comment out languages you don't need.
# OR 
# Add new languages
# RUN apt install -y python3  --- python3 is already installed in bullseye distro
# RUN apt install -y python
# RUN apt install -y nodejs
# RUN apt install -y openjdk-17-jdk 

COPY . .

CMD ["go","run","-mod=vendor","main.go"]
