package config

import "testing"

func TestLoadDefaultRAGConfig(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}

	if !cfg.RAG.Enabled {
		t.Fatal("rag should be enabled by default")
	}
	if cfg.RAG.Provider != "local" {
		t.Fatalf("unexpected rag provider: %s", cfg.RAG.Provider)
	}
	if cfg.RAG.Address != "115.190.209.83:19530" {
		t.Fatalf("unexpected rag address: %s", cfg.RAG.Address)
	}
	if cfg.RAG.CollectionName != "travel_knowledge" {
		t.Fatalf("unexpected rag collection: %s", cfg.RAG.CollectionName)
	}
	if cfg.RAG.TopK != 3 {
		t.Fatalf("unexpected rag top k: %d", cfg.RAG.TopK)
	}
	if cfg.RAG.MinScore <= 0 {
		t.Fatalf("min score should be positive: %f", cfg.RAG.MinScore)
	}
}

func TestLoadDefaultInfrastructureConfig(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}

	if cfg.MySQL.DSN == "" {
		t.Fatal("mysql dsn should be configured")
	}
	if cfg.Redis.Addr == "" {
		t.Fatal("redis addr should be configured")
	}
	if cfg.Auth.JWTSecret == "" {
		t.Fatal("jwt secret should be configured")
	}
	if cfg.Auth.JWTExpire <= 0 {
		t.Fatalf("jwt expire should be positive: %s", cfg.Auth.JWTExpire)
	}
}
