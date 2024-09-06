package controllers

type ApiConfig struct {
	fileserverHits int
	jwtSecret      string
}

func NewApiConfig(secret string) ApiConfig {
	return ApiConfig{
		fileserverHits: 0,
		jwtSecret:      secret,
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
