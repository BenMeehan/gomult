
# GOMULT

Gomult is a bunch of compilers for the most popular languages in a dockerfile, providing a gRPC interface and written in golang. 


## How it works

All the compilers/interpreters needed are installed using the 'RUN' directive in Dockerfile.

The Golang gRPC server then listens on port 9000 and any program can be compiled by passing the language and code using a gRPC client in any of the [supported](https://grpc.io/docs/languages/) languages.

For a sample gRPC client written in GO, checkout the test directory in the repo.

## Supported languages

This project is in early stage and currently supports

- Python 3
- Python 2.7
- NodeJS 12
- Java 17 (code must have atleast 1 public class)
- Golang 1.19.3
- C 

with support for more languages to be added soon. 

## NOTE

Installing all the Compilers/Interpreters may be heavy and use up disk space. So if you are using this project either

- Make sure you have enough disk space
- Comment out the Compilers/Interpreters for languages you don't need in the dockerfile and build it.

## Authors

- [@benmeehan111](https://github.com/BenMeehan)


## Acknowledgements

 - [gRPC](https://github.com/grpc/grpc)
 - [docker](https://www.docker.com/)
 - [Debian "bullseye"](https://www.debian.org/releases/bullseye/)

