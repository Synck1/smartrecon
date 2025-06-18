package smartrecon

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"smartrecon/config"
	"smartrecon/core"
	"strings"
)

func setupOutputDir() {
	os.MkdirAll("output", 0755)
}

func askDomain() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Digite o domínio (ex: globo.com.br): ")
	domainInput, _ := reader.ReadString('\n')
	return strings.TrimSpace(domainInput)
}

func getFromCrtSh(domain string) []string {
	fmt.Println("[*] Coletando subdomínios via crt.sh...")
	crtshSubs, err := core.FetchFromCrtSh(domain)
	if err != nil {
		fmt.Println("[-] Erro no crt.sh:", err)
		return []string{}
	}
	fmt.Printf("[+] %d subdomínios via crt.sh\n", len(crtshSubs))
	return crtshSubs
}

func getFromRevWhois(domain string) []string {
	fmt.Println("[*] Buscando domínios relacionados via RevWhois...")
	relatedDomains, err := core.RunRevWhois(domain)
	if err != nil {
		fmt.Println("[-] Erro no RevWhois:", err)
		return []string{}
	}
	fmt.Printf("[+] %d domínios relacionados encontrados via RevWhois\n", len(relatedDomains))
	return relatedDomains
}

func getFromAmass(domain string) []string {
	fmt.Println("[*] Coletando subdomínios via Amass...")
	subs, err := core.RunAmass(domain)
	if err != nil {
		fmt.Println("[-] Erro no Amass:", err)
		return []string{}
	}
	fmt.Printf("[+] %d subdomínios via Amass\n", len(subs))
	return subs
}

func getFromSubfinder(domain string) []string {
	fmt.Println("[*] Coletando subdomínios via Subfinder...")
	subs, err := core.RunSubfinder(domain)
	if err != nil {
		fmt.Println("[-] Erro no Subfinder:", err)
		return []string{}
	}
	fmt.Printf("[+] %d subdomínios via Subfinder\n", len(subs))
	return subs
}

func generatePermutations(subs []string) []string {
	fmt.Println("[*] Gerando permutações inteligentes...")
	perms := core.GenerateAutoPermutations(subs)
	fmt.Printf("[+] %d subdomínios gerados por permutação\n", len(perms))
	return perms
}

func saveSubdomains(subs []string) {
	if err := core.SaveToFile(subs, "output/subs.txt"); err != nil {
		log.Println("[-] Erro ao salvar subs.txt:", err)
	} else {
		fmt.Println("[+] Subdomínios salvos em output/subs.txt")
	}
}

func resolveDNS() []string {
	fmt.Println("[*] Rodando DNSX para resolver subdomínios...")
	err := core.RunDNSX("output/subs.txt", "output/resolved.txt")
	if err != nil {
		log.Println("[-] Erro ao rodar dnsx:", err)
		return []string{}
	}
	fmt.Println("[+] Subdomínios resolvidos salvos em output/resolved.txt")
	ips := core.ExtractIPsFromDNSXOutput("output/resolved.txt")
	core.SaveToFile(ips, "output/ips.txt")
	fmt.Printf("[+] Total de %d IPs resolvidos\n", len(ips))
	return ips
}

func runASNMapping(ips []string) {
	fmt.Println("[*] Mapeando IPs → ASN → Ranges...")

	var allRanges []string

	for _, ip := range ips {
		asn, owner, cidr := core.IPToASN(ip)
		if asn != "" {
			fmt.Printf("[ASN INFO] IP: %s | ASN: %s | OWNER: %s | CIDR: %s\n", ip, asn, owner, cidr)

			ranges, err := core.ASNToRanges(asn)
			if err == nil {
				allRanges = append(allRanges, ranges...)
			} else {
				fmt.Println("[-] Erro ao obter ranges do ASN:", err)
			}
		}
	}

	allRanges = core.CleanLines(allRanges)
	core.SaveToFile(allRanges, "output/ranges.txt")
	fmt.Printf("[+] Total de %d ranges salvos em output/ranges.txt\n", len(allRanges))
}

func asnFilter() {

	keywords := []string{"globo", "globosat", "grupo globo"}

	asns := core.FilterASNFromIPs("output/ips.txt", keywords)

	fmt.Println("[+] ASN encontrados que batem com a empresa:")
	for _, a := range asns {
		fmt.Println(a)
	}
}

func Run(cfg *config.Config) {
	setupOutputDir()

	var domain string
	if cfg.Domain != "" {
		domain = cfg.Domain
		fmt.Println("[*] Usando domínio do config:", domain)
	} else {
		domain = askDomain()
	}

	var allSubs []string

	allSubs = append(allSubs, getFromCrtSh(domain)...)
	allSubs = append(allSubs, getFromSubfinder(domain)...)
	allSubs = append(allSubs, getFromAmass(domain)...)

	relatedDomains := getFromRevWhois(domain)

	for _, related := range relatedDomains {
		fmt.Println("[*] Coletando subdomínios do domínio relacionado:", related)
		allSubs = append(allSubs, getFromCrtSh(related)...)
		allSubs = append(allSubs, getFromSubfinder(related)...)
		allSubs = append(allSubs, getFromAmass(related)...)
	}

	allSubs = core.CleanLines(allSubs)
	fmt.Printf("[+] %d subdomínios únicos coletados\n", len(allSubs))

	perms := generatePermutations(allSubs)
	allFinal := core.CleanLines(append(allSubs, perms...))
	fmt.Printf("[+] %d subdomínios após permutação\n", len(allFinal))

	saveSubdomains(allFinal)

	ips := resolveDNS()

	runASNMapping(ips)
	asnFilter()

	fmt.Println("[*] Recon Finalizado com sucesso!")
}
