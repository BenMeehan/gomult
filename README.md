
# GOMULT

Gomult is a multi-language compiler built into a dockerfile, leveraging the power of gRPC and speed of Golang.


## How it works

All the compilers/interpreters needed are installed using the 'RUN' directive in Dockerfile.

The Golang gRPC server then listens on port 9000 and any program can be compiled by passing the language and code using a gRPC client in any of the [supported](https://grpc.io/docs/languages/) languages.

For a sample gRPC client written in GO, checkout the test directory in the repo.
## Supported languages

This project is in early stage and currently supports

- Python 3
- Python 2.7
- Java 17
- Golang

with support for more languages to be added soon. 

You can even add your own language. Just follow these steps,

1. Fork or clone the repo
2. Edit the Dockerfile to add instructions to install your compiler (or comment out the ones you don't need).
3. Add go code in the compiler.go and helpers.go to handle the compilation.
4. Build the dockerfile.
## Authors

- [@benmeehan111](https://github.com/BenMeehan)


## Acknowledgements

 - [gRPC](https://github.com/grpc/grpc)
 - [docker](https://www.docker.com/)
 - [Debian "bullseye"](https://www.debian.org/releases/bullseye/)

