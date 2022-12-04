package code

import (
	"log"
	"fmt"
	"os"
	"os/exec"
	"github.com/google/uuid"
)

func handleError(err error) {
	if err != nil {
		log.Println("Error: ",err)
	}
}

func compilePython3(code string)string {
	id := uuid.New().String()
	file := fmt.Sprintf("/tmp/%s.py", id)
	f, err := os.Create(file)
	handleError(err)
	_, err = f.WriteString(code)
	handleError(err)
	out, err := exec.Command("python3", file).CombinedOutput()
	handleError(err)
	err = os.Remove(file)
	handleError(err)
	return string(out[:])
}
