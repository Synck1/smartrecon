package main

import (
	"fmt"
	"log"
	"os"
	"smartrecon/core"
)

func main() {

	os.MkdirAll("output", 0755)

	domain := "example.com"

	var allSubs []string

	// 🔎 RevWHOIS
	fmt.Println("[*] Buscando domínios relacionados via RevWhois...")
	relatedDomains, err := core.RunRevWhois(domain)
	if err == nil {
		fmt.Printf("[+] %d domínios relacionados encontrados\n", len(relatedDomains))
	} else {
		fmt.Println("[-] Erro no RevWhois:", err)
	}

	// 🔍 CRT.sh
	fmt.Println("[*] Coletando subdomínios via crt.sh...")
	crtshSubs, err := core.FetchFromCrtSh(domain)
	if err == nil {
		fmt.Printf("[+] %d subdomínios via crt.sh\n", len(crtshSubs))
		allSubs = append(allSubs, crtshSubs...)
	} else {
		fmt.Println("[-] Erro no crt.sh:", err)
	}

	// 🔍 Subfinder
	fmt.Println("[*] Coletando com subfinder...")
	subs1, err := core.RunSubfinder(domain)
	if err == nil {
		allSubs = append(allSubs, subs1...)
	}

	// 🔍 Amass
	fmt.Println("[*] Coletando com amass...")
	// subs2, err := core.RunAmass(domain)
	// if err == nil {
	// 	allSubs = append(allSubs, subs2...)
	// }

	// 🔁 Junta e limpa
	allSubs = core.CleanLines(allSubs)
	fmt.Printf("[+] Total de %d subdomínios únicos coletados\n", len(allSubs))

	// 🤖 Permutações inteligentes
	fmt.Println("[*] Gerando permutações inteligentes...")
	perms := core.GenerateAutoPermutations(allSubs)
	allFinal := core.CleanLines(append(allSubs, perms...))
	fmt.Printf("[+] %d subdomínios após permutação\n", len(allFinal))

	// 💾 Salva
	if err := core.SaveToFile(allFinal, "output/subs.txt"); err != nil {
		log.Println("[-] Erro ao salvar subs.txt:", err)
	}

	// Resolução com shuffledns
	/*
		fmt.Println("[*] Resolvendo com shuffledns...")
		resolved, err := core.RunShuffledns(domain, "data/dns.txt", "data/resolvers.txt")
		if err == nil {
			fmt.Printf("[+] %d ativos\n", len(resolved))
			core.SaveToFile(resolved, "output/resolved.txt")
		}
	*/
}
