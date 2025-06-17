package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// 🔥 Struct que mapeia seu config.yaml
type Config struct {
	Domain string `yaml:"domain"`

	Modules struct {
		CRTSh         bool `yaml:"crtsh"`
		Subfinder     bool `yaml:"subfinder"`
		Permutations  bool `yaml:"permutations"`
		RevWhois      bool `yaml:"revwhois"`
		Amass         bool `yaml:"amass"`
		DNSResolution bool `yaml:"dns_resolution"`
		HTTPCheck     bool `yaml:"http_check"`
	} `yaml:"modules"`

	Permutations struct {
		Wordlist string `yaml:"wordlist"`
		MaxDepth int    `yaml:"max_depth"`
	} `yaml:"permutations"`
}

// 🚀 Função que carrega o arquivo YAML pra dentro da struct
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	// Limpa espaços no domínio (boa prática)
	cfg.Domain = strings.TrimSpace(cfg.Domain)

	return &cfg, nil
}
