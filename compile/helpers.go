package compile

import (
	fmt "fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func handleError(err error) {
	if err != nil {
		log.Error(err)
	}
}

func getMainClass(code string) string {
	file, err := ioutil.TempFile("/tmp", "java")
	handleError(err)
	_, err = file.WriteString(code)
	handleError(err)
	defer os.Remove(file.Name())
	out, err := exec.Command("bash", "-c", fmt.Sprintf(`awk '/public class/' %s;`, file.Name())).CombinedOutput()
	handleError(err)
	name := string(out[13:])
	end := strings.Index(name, "{")
	return name[:end]
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
	fname := fmt.Sprintf("/tmp/%s.java", mainClass)
	f, err := os.Create(fname)
	handleError(err)
	_, err = f.WriteString(code)
	handleError(err)
	_, err = exec.Command("javac", fname).Output()
	handleError(err)
	out, err := exec.Command("java", "-cp", "/tmp", mainClass).CombinedOutput()
	handleError(err)
	err = os.Remove(fname)
	handleError(err)
	return string(out[:])
}
