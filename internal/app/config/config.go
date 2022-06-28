package config

import (
	"flag"
	"phrasetagg/url-shortener/internal/app/helpers"
)

type cfg struct {
	FileStoragePath string
	BaseURL         string
	ServerAddr      string
}

func PrepareCfg() *cfg {
	cfg := new(cfg)

	flag.StringVar(&cfg.FileStoragePath, "f", helpers.GetEnv("FILE_STORAGE_PATH", ""), "URLs file path")
	flag.StringVar(&cfg.BaseURL, "b", helpers.GetEnv("BASE_URL", "http://localhost:8080/"), "short URLs base URL")
	flag.StringVar(&cfg.ServerAddr, "a", helpers.GetEnv("SERVER_ADDRESS", "localhost:8080"), "host:port of the server")

	flag.Parse()

	lastChar := cfg.BaseURL[len(cfg.BaseURL)-1:]
	if lastChar != "/" {
		cfg.BaseURL = cfg.BaseURL + "/"
	}

	return cfg
}
