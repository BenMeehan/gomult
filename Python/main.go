package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// Maximum allowed code size in bytes
const MaxCodeSize = 1024 * 1024 // 1 MB

// Restricted user and group ID - always 1000 for docker unless explicit
const RestrictedUserID = 1000
const RestrictedGroupID = 1000

type CompileRequest struct {
	Code     string `json:"code"`
	Input    string `json:"input"`
	Language string `json:"language,omitempty"`
}

func main() {
	http.HandleFunc("/compile", handleCompile)
	fmt.Println("Python Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCompile(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error occurred during parsing of the request body", err)
		log.Printf("Failed to read request body: %v", err)
		return
	}

	// Parse the JSON request body
	var compileReq CompileRequest
	err = json.Unmarshal(body, &compileReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error occurred during parsing of the JSON request body", err)
		log.Printf("Failed to parse JSON request body: %v", err)
		return
	}

	code := []byte(compileReq.Code)

	// Check if the code size exceeds the maximum allowed size
	if len(code) > MaxCodeSize {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Code size too big", err)
		log.Printf("Code size exceeds the maximum limit of %d bytes", MaxCodeSize)
		return
	}

	// Create a temporary file to store the code
	tmpFile, err := ioutil.TempFile("", "code-*.py")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal server error", err)
		log.Printf("Failed to create temporary file: %v", err)
		return
	}
	// defer os.Remove(tmpFile.Name()) // Clean up the temporary file

	// Set the desired permissions for the temporary file
	err = os.Chmod(tmpFile.Name(), 0050)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal server error", err)
		log.Printf("Failed to set permissions for temporary file: %v", err)
		return
	}

	// Write the code to the temporary file
	_, err = tmpFile.Write(code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal server error", err)
		log.Printf("Failed to write code to temporary file: %v", err)
		return
	}

	// Close the temporary file
	err = tmpFile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal server error", err)
		log.Printf("Failed to close temporary file: %v", err)
		return
	}

	// Create a context with a timeout duration
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a channel to receive the output
	outputChannel := make(chan []byte)

	// Run the Python code in a Goroutine and monitor for timeouts
	go func() {
		cmd := exec.CommandContext(ctx, "python3", tmpFile.Name())

		// Set the user and group ID of the executed program
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: RestrictedUserID,
				Gid: RestrictedGroupID,
			},
		}

		cmdOutput, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Execution error: %s", err)
		}

		// Send the execution output through the channel
		outputChannel <- cmdOutput
	}()

	select {
	case <-ctx.Done():
		// Execution timed out
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Execution timed out")
		log.Println("Execution timed out")
	case output := <-outputChannel:
		// Execution completed within the timeout duration
		w.Header().Set("Content-Type", "text/plain")
		w.Write(output)
	}
}
