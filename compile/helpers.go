package compile

import (
	"os/exec"

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
