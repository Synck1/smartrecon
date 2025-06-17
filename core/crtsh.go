package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
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
