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

// Cria pasta output se não existir
func setupOutputDir() {
	os.MkdirAll("output", 0755)
}

// Pergunta o domínio pro usuário
func askDomain() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Digite o domínio (ex: globo.com.br): ")
	domainInput, _ := reader.ReadString('\n')
	return strings.TrimSpace(domainInput)
}

// Busca no CRT.sh
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

// Busca no RevWhois (domínios relacionados)
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

// Busca no Amass
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

// Busca no Subfinder
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

// Gera permutações
func generatePermutations(subs []string) []string {
	fmt.Println("[*] Gerando permutações inteligentes...")
	perms := core.GenerateAutoPermutations(subs)
	fmt.Printf("[+] %d subdomínios gerados por permutação\n", len(perms))
	return perms
}

// Salva resultado
func saveSubdomains(subs []string) {
	if err := core.SaveToFile(subs, "output/subs.txt"); err != nil {
		log.Println("[-] Erro ao salvar subs.txt:", err)
	} else {
		fmt.Println("[+] Subdomínios salvos em output/subs.txt")
	}
}

// Resolve DNS
func resolveDNS() {
	fmt.Println("[*] Rodando DNSX para resolver subdomínios...")
	if err := core.RunDNSX("output/subs.txt", "output/resolved.txt"); err != nil {
		log.Println("[-] Erro ao rodar dnsx:", err)
	} else {
		fmt.Println("[+] Subdomínios resolvidos salvos em output/resolved.txt")
	}
}

// Função principal
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

	// 🔥 Busca subdomínios do domínio principal
	allSubs = append(allSubs, getFromCrtSh(domain)...)
	allSubs = append(allSubs, getFromSubfinder(domain)...)
	allSubs = append(allSubs, getFromAmass(domain)...)

	// 🔍 Busca domínios relacionados
	relatedDomains := getFromRevWhois(domain)

	// 🔥 Roda o mesmo recon para os domínios relacionados
	for _, related := range relatedDomains {
		fmt.Println("[*] Coletando subdomínios do domínio relacionado:", related)
		allSubs = append(allSubs, getFromCrtSh(related)...)
		allSubs = append(allSubs, getFromSubfinder(related)...)
		allSubs = append(allSubs, getFromAmass(related)...)
	}

	// 🔁 Limpa duplicatas
	allSubs = core.CleanLines(allSubs)
	fmt.Printf("[+] %d subdomínios únicos coletados\n", len(allSubs))

	// 🤖 Permutação inteligente
	perms := generatePermutations(allSubs)
	allFinal := core.CleanLines(append(allSubs, perms...))
	fmt.Printf("[+] %d subdomínios após permutação\n", len(allFinal))

	// 💾 Salvar e Resolver
	saveSubdomains(allFinal)
	resolveDNS()

	fmt.Println("[*] Recon Finalizado com sucesso!")
}
