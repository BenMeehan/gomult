package compile

import (
	"os"
	"os/exec"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func handleError(err error) {
	if err != nil {
		log.Error(err)
	}
}

func compilePython3(code string) string {
	out, err := exec.Command("python3", "-c", code).CombinedOutput()
	handleError(err)
	return string(out[:])
}

func compilePython(code string) string {
	out, err := exec.Command("python", "-c", code).CombinedOutput()
	handleError(err)
	return string(out[:])
}

func compileNodeJS(code string) string {
	out, err := exec.Command("node", "-e", code).CombinedOutput()
	handleError(err)
	return string(out[:])
}

func compileJava17(code string, mainClass string) string {
	fname := createFile(mainClass, "java", code)
	_, err := exec.Command("javac", fname).Output()
	handleError(err)
	out, err := exec.Command("java", "-cp", "/tmp", mainClass).CombinedOutput()
	handleError(err)
	err = os.Remove(fname)
	handleError(err)
	return string(out[:])
}

func compileGolang(code string) string {
	id := uuid.New().String()
	fname := createFile(id, "go", code)
	out, err := exec.Command("go", "run", fname).CombinedOutput()
	handleError(err)
	err = os.Remove(fname)
	handleError(err)
	return string(out[:])
}
