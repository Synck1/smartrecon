package core

import (
	"os/exec"
	"strings"
)

func RunAmass(domain string) ([]string, error) {
	cmd := exec.Command("amass", "enum", "-passive", "-d", domain, "-silent")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	return lines, nil
}
