package core

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

func SaveToFile(lines []string, path string) error {
	data := strings.Join(lines, "\n")
	return ioutil.WriteFile(path, []byte(data), 0644)
}

func CleanLines(lines []string) []string {
	seen := make(map[string]bool)
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !seen[line] {
			seen[line] = true
			cleaned = append(cleaned, line)
		}
	}
	return cleaned
}

func LoadWordlist(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	return words, scanner.Err()
}
