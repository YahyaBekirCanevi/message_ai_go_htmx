package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	geminiAPIKey   string
	geminiAPIURL   string
	loadConfigOnce sync.Once
)

type AppConfig struct {
	GeminiAPIKey string `yaml:"gemini_api_key"`
	GeminiAPIURL string `yaml:"gemini_api_url"`
}

func LoadConfig() (string, string) {
	loadConfigOnce.Do(func() {
		f, err := os.Open("application.yml")
		if err != nil {
			fmt.Println("Warning: could not open application.yml:", err)
			return
		}
		defer f.Close()
		var cfg AppConfig
		if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
			fmt.Println("Warning: could not decode application.yml:", err)
			return
		}
		geminiAPIKey = cfg.GeminiAPIKey
		geminiAPIURL = cfg.GeminiAPIURL
	})
	return geminiAPIKey, geminiAPIURL
}
