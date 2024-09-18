package controllers

type ApiConfig struct {
	fileserverHits int
	jwtSecret      string
	polkaApiKey    string
}

func NewApiConfig(jwtSecret, polkaApiKey string) ApiConfig {
	return ApiConfig{
		fileserverHits: 0,
		jwtSecret:      jwtSecret,
		polkaApiKey:    polkaApiKey,
	}
}

func (cfg *ApiConfig) ServerHitsReset() {
	cfg.fileserverHits = 0
}

func (cfg *ApiConfig) ServerHitsIncrement() {
	cfg.fileserverHits++
}

func (cfg *ApiConfig) ServerHitsGet() int {
	return cfg.fileserverHits
}
