package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

// Maximum allowed code size in bytes
const MaxCodeSize = 1024 * 1024 // 1 MB

// Restricted user and group ID
const RestrictedUserID = 1000
const RestrictedGroupID = 1000

func main() {
	http.HandleFunc("/compile", handleCompile)
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCompile(w http.ResponseWriter, r *http.Request) {
	// Read the code from the HTTP request body
	code, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Failed to read code from request body: %v", err)
		return
	}

	// Check if the code size exceeds the maximum allowed size
	if len(code) > MaxCodeSize {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Code size exceeds the maximum limit of %d bytes", MaxCodeSize)
		return
	}

	// Create a temporary file to store the code
	tmpFile, err := ioutil.TempFile("", "code-*.c")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to create temporary file: %v", err)
		return
	}
	defer os.Remove(tmpFile.Name()) // Clean up the temporary file

	// Write the code to the temporary file
	_, err = tmpFile.Write(code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to write code to temporary file: %v", err)
		return
	}

	// Close the temporary file
	err = tmpFile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to close temporary file: %v", err)
		return
	}

	// Compile the code using GCC (assuming GCC is installed)
	outputFile := filepath.Join(filepath.Dir(tmpFile.Name()), "output")
	cmd := exec.Command("gcc", tmpFile.Name(), "-o", outputFile)

	// Create a new syscall.SysProcAttr and set its Ptrace to false
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: false,
	}

	compilerOutput, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Compilation error: %s", compilerOutput)
		log.Printf("Compilation error: %s", err)
		return
	}

	log.Printf("Compilation successful. Output file: %s", outputFile)

	// Create a context with a timeout duration
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a channel to receive the output
	outputChannel := make(chan []byte)

	// Run the compiled output in a Goroutine and monitor for timeouts
	go func() {
		cmd := exec.CommandContext(ctx, outputFile)

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

