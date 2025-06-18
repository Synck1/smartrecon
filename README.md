# Synck SmartRecon


![Logo do projeto](/static/synck-logo.png.png)


> Ferramenta de Reconhecimento (Recon) desenvolvida em Go para coleta de informações de domínios e subdomínios.

---

[![Go version](https://img.shields.io/badge/go-1.18+-00ADD8?logo=go&style=for-the-badge)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/rafasor/smartrecon?style=for-the-badge)](https://github.com/rafasor/smartrecon/issues)
[![GitHub stars](https://img.shields.io/github/stars/rafasor/smartrecon?style=for-the-badge)](https://github.com/rafasor/smartrecon/stargazers)

---

## 📋 Índice

- [Sobre](#-sobre)
- [Funcionalidades](#-funcionalidades)
- [Instalação](#-instalacao)
- [Uso](#-uso)
- [Configuração](#-configuracao)
- [Contribuição](#-contribuicao)
- [Licença](#-licenca)
- [Contato](#-contato)

---

## 🔎 Sobre

SmartRecon é uma ferramenta rápida e extensível escrita em Go, focada em realizar reconhecimento e enumeração de domínios, subdomínios e informações relacionadas para apoiar pentests e pesquisas em segurança.

---

## ⚙️ Funcionalidades

- Enumeração de subdomínios via múltiplas fontes
- Resolução DNS integrada
- Checagem de status HTTP
- Permutações e variações de nomes de subdomínios
- Suporte a módulos configuráveis
- Relatórios simples e exportação em JSON ( Em Desenvolvimento )

---

## 🚀 Instalação

### Pré-requisitos

- Go 1.18+ instalado

### Clonando e compilando

```bash
git clone https://github.com/rafasor/smartrecon.git
cd smartrecon
go build -o smartrecon main.go
```

---

## 🛠 Uso

Executando a ferramenta básica:

```bash
./smartrecon -domain example.com
```

Exemplo com opções avançadas:

```bash
./smartrecon -domain example.com -modules crtsh,amass,subfinder -output results.json
```

---

## ⚙️ Configuração

Edite o arquivo `config.yaml` para ajustar módulos ativados, listas de permutações e outras opções.

---

## 🤝 Contribuição

Contribuições são bem-vindas! Por favor, abra issues ou pull requests para melhorias, correções e novas funcionalidades.

---

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para detalhes.

---

## 📬 Contato

Link do projeto: [https://github.com/rafasor/smartrecon](https://github.com/rafasor/smartrecon)
