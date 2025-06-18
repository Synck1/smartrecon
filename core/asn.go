package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// 🔥 Remove ANSI (se quiser aplicar na leitura de DNSX)
func RemoveANSICodes(input string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(input, "")
}

// 🔍 Extrai IPs de arquivo DNSX
func ExtractIPsFromDNSXOutput(path string) []string {
	data, _ := os.ReadFile(path)
	lines := strings.Split(string(data), "\n")
	var ips []string

	re := regexp.MustCompile(`\[(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\]`)

	for _, line := range lines {
		line = RemoveANSICodes(line)
		match := re.FindStringSubmatch(line)
		if len(match) > 1 {
			ips = append(ips, match[1])
		}
	}

	return CleanLines(ips)
}

// 🔗 IP → ASN → CIDR
func IPToASN(ip string) (asn string, owner string, cidr string) {
	url := fmt.Sprintf("https://ipinfo.io/%s/json", ip)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("[-] Erro IPToASN:", err)
		return "", "", ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	asn = fmt.Sprintf("%v", result["org"])
	cidr = fmt.Sprintf("%v", result["route"])
	owner = asn

	return asn, owner, cidr
}

// ASN → Todos os Ranges
func ASNToRanges(asn string) ([]string, error) {
	url := fmt.Sprintf("https://api.bgpview.io/asn/%s/prefixes", strings.TrimPrefix(asn, "AS"))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Data struct {
			IPv4Prefixes []struct {
				Prefix string `json:"prefix"`
			} `json:"ipv4_prefixes"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	var ranges []string
	for _, p := range result.Data.IPv4Prefixes {
		ranges = append(ranges, p.Prefix)
	}

	return ranges, nil
}

// 🔥 🔥 Filtra ASN que pertence à empresa
func FilterASNFromIPs(ipListPath string, keywords []string) []string {
	file, err := os.Open(ipListPath)
	if err != nil {
		fmt.Println("[-] Erro ao abrir arquivo:", err)
		return []string{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var asnList []string

	for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())
		if ip == "" {
			continue
		}

		asn, owner, cidr := IPToASN(ip)
		if asn == "" {
			continue
		}

		fmt.Printf("[ASN INFO] IP: %s | ASN: %s | OWNER: %s | CIDR: %s\n", ip, asn, owner, cidr)

		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(owner), strings.ToLower(keyword)) {
				entry := fmt.Sprintf("%s | %s | %s", asn, owner, cidr)
				asnList = append(asnList, entry)
			}
		}
	}

	// Salvar resultado
	SaveToFile(asnList, "output/asn_empresa.txt")
	fmt.Println("[+] Lista de ASN da empresa salva em output/asn_empresa.txt")

	return CleanLines(asnList)
}
