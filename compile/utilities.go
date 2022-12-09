package compile

import (
	fmt "fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func createFile(name string, extension string, contents string) string {
	fname := fmt.Sprintf("/tmp/%s.%s", name, extension)
	f, err := os.Create(fname)
	handleError(err)
	_, err = f.WriteString(contents)
	handleError(err)
	return fname
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
