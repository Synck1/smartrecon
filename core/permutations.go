package core

import (
	"strings"
)

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

		name := parts[0]
		domain := strings.Join(parts[1:], ".")

		for _, word := range wordlist {
			word = strings.TrimSpace(word)
			if word == "" {
				continue
			}

			mutated = append(mutated, word+"."+sub)
			mutated = append(mutated, name+"-"+word+"."+domain)
			mutated = append(mutated, word+"-"+name+"."+domain)
			mutated = append(mutated, name+"."+word+"."+domain)
		}
	}
	return mutated
}

func GenerateAutoPermutations(subdomains []string) []string {
	wordlist := ExtractWords(subdomains)
	return GeneratePermutations(subdomains, wordlist)
}

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
