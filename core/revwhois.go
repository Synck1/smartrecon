package core

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

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
