package core

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func FetchFromCrtSh(domain string) ([]string, error) {
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	re := regexp.MustCompile(`\"name_value\":\"([^\"]+)\"`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	found := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			for _, entry := range strings.Split(match[1], "\\n") {
				if strings.HasSuffix(entry, domain) {
					found[entry] = true
				}
			}
		}
	}

	var results []string
	for sub := range found {
		results = append(results, sub)
	}
	return results, nil
}

// Subfinder wrapper
func RunSubfinder(domain string) ([]string, error) {
	cmd := exec.Command("subfinder", "-d", domain, "-silent")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	return lines, nil
}

// Amass wrapper
// func RunAmass(domain string) ([]string, error) {
// 	cmd := exec.Command("amass", "enum", "-passive", "-d", domain)
// 	out, err := cmd.Output()
// 	if err != nil {
// 		return nil, err
// 	}
// 	lines := strings.Split(string(out), "\n")
// 	return lines, nil
// }

// Shuffledns wrapper
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

// RevWhois using viewdns.info
func RunRevWhois(query string) ([]string, error) {
	url := "https://viewdns.info/reversewhois/?q=" + query
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	re := regexp.MustCompile(`(?i)<td>([a-z0-9.-]+\.[a-z]{2,})</td>`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	var domains []string
	for _, m := range matches {
		if len(m) > 1 {
			domains = append(domains, m[1])
		}
	}
	return domains, nil
}

// Utilitário: salva em arquivo
func SaveToFile(lines []string, path string) error {
	data := strings.Join(lines, "\n")
	return ioutil.WriteFile(path, []byte(data), 0644)
}

// Utilitário: remove duplicatas e entradas vazias
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

// Carrega wordlist de arquivo
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

// Extrai palavras únicas de subdomínios já encontrados
func ExtractWords(subdomains []string) []string {
	seen := make(map[string]bool)
	var words []string

	for _, sub := range subdomains {
		sub = strings.TrimSpace(sub)
		if sub == "" {
			continue
		}
		parts := strings.Split(sub, ".")
		for _, part := range parts {
			if part == "com" || part == "www" || part == "net" || part == "org" {
				continue
			}
			if len(part) < 2 || len(part) > 20 {
				continue
			}
			if !seen[part] {
				seen[part] = true
				words = append(words, part)
			}
		}
	}
	return words
}

// Gera permutações usando uma wordlist fornecida
func GeneratePermutations(subdomains []string, wordlist []string) []string {
	var mutated []string
	for _, sub := range subdomains {
		sub = strings.TrimSpace(sub)
		if sub == "" {
			continue
		}

		parts := strings.Split(sub, ".")
		if len(parts) < 3 {
			continue
		}

		name := parts[0]                       // ex: "api"
		domain := strings.Join(parts[1:], ".") // ex: "example.com"

		for _, word := range wordlist {
			word = strings.TrimSpace(word)
			if word == "" {
				continue
			}

			mutated = append(mutated, word+"."+sub)             // dev.api.example.com
			mutated = append(mutated, name+"-"+word+"."+domain) // api-dev.example.com
			mutated = append(mutated, word+"-"+name+"."+domain) // dev-api.example.com
			mutated = append(mutated, name+"."+word+"."+domain) // api.dev.example.com
		}
	}
	return mutated
}

// Gera permutações usando wordlist dinâmica baseada nos subdomínios
func GenerateAutoPermutations(subdomains []string) []string {
	wordlist := ExtractWords(subdomains)
	return GeneratePermutations(subdomains, wordlist)
}
