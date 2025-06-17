package main

import (
	"fmt"
	"log"
	smartrecon "smartrecon/cmd"
	"smartrecon/config"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatal("Erro carregando configuração:", err)
	}

	fmt.Println("[*] SmartRecon iniciado...")
	smartrecon.Run(cfg)
}
