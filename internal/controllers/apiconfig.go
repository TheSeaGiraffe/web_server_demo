package controllers

type ApiConfig struct {
	fileserverHits int
}

func NewApiConfig() ApiConfig {
	return ApiConfig{
		fileserverHits: 0,
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
