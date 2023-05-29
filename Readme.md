# Code Compilation and Execution Service

This project provides a compilation and execution service for multiple programming languages. It allows you to submit code snippets in various languages and receive the corresponding output.

### Features

- Supports multiple programming languages
- Secure execution environment with restricted user privileges
- Handles code compilation and execution
- Timeout mechanism to prevent long-running executions


### Prerequisites
- Docker: Make sure you have Docker installed on your machine to run the code execution service.

### Getting Started

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

#### Request Format:
```json
{
  "code": "<code_snippet>",
  "input": "<input_for_the_program>",
  "language": "<programming_language>"
}
```