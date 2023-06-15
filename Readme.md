# Code Compilation and Execution Service

This project provides a compilation and execution service for multiple programming languages. It allows you to submit code snippets in various languages and receive the corresponding output.

For demo API's please look at [API Readme](./Api.md)

Thank you [Render](https://render.com/) for your free services.

## Features

- Supports multiple programming languages
- Secure execution environment with restricted user privileges
- Handles code compilation and execution
- Timeout mechanism to prevent long-running executions


## Prerequisites
- Docker: Make sure you have Docker installed on your machine to run the code execution service.

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/your-username/your-repo.git
```

2. Navigate to the project directory:
```bash
cd your-repo
```

3. Build the Docker image:
```bash
docker build -t code-execution-service .
```
4. Run the Docker container:
```bash
docker run -p 8080:8080 code-execution-service
```

5. The code execution service is now running. You can access it at `http://localhost:8080`.


## API Usage
### Endpoint: /compile

Send a POST request to this endpoint to compile and execute code.

### Request Format:
```json
{
  "code": "<code_snippet>",
  "input": "<input_for_the_program>",
  "language": "<programming_language>"
}
```

- **code:** The code snippet to be compiled and executed.
- **input:** The input to be provided to the program (optional, depending on the language).
- **language:** The programming language of the code snippet.


**Make a post request to this URL(https://load-balancer-1l8h.onrender.com) to test it out.**

*Example curl request*

```
curl -X POST -H "Content-Type: application/json" -d '{
  "code": "print(\"Hello, World!\")",
  "input": "",
  "language": "py"
}' https://load-balancer-1l8h.onrender.com/compile
```

### Response Format:

If the compilation and execution are successful within the timeout duration, the API will respond with the output of the program.

If an error occurs during compilation or execution, the API will respond with an appropriate error message.

## Supported Languages

The service currently supports the following programming languages:

- C
- C++
- Python 3
- Python 2.7
- Java 11
- Javascript
- Golang
- Rust
- more to be added soon...

You can extend the service to support additional languages by adding the corresponding code compilation and execution logic.

## Contributing

Contributions are welcome! If you would like to contribute to this project, please follow these steps:

- Fork the repository.
- Create a new branch for your feature or bug fix.
- Make your modifications.
- Commit your changes and push the branch to your forked repository.
- Submit a pull request detailing your changes.

## License
This project is licensed under the [Apache License](./LICENSE)

## Contact 
For any inquiries or support, please contact benmeehan111@gmail.com

