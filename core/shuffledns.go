package core

import (
	"os/exec"
	"strings"
)

func RunShuffledns(domain, wordlistPath, resolversPath string) ([]string, error) {
	cmd := exec.Command("shuffledns",
		"-d", domain,
		"-w", wordlistPath,
		"-r", resolversPath,
		"-silent",
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	return lines, nil
}
