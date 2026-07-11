package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

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
	if cfg.RAG.Address != "127.0.0.1:19530" {
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
	if cfg.SMS.Endpoint == "" {
		t.Fatal("sms endpoint should be configured")
	}
	if cfg.SMS.RegisterCodeExpire != 5*time.Minute {
		t.Fatalf("unexpected register code expire: %s", cfg.SMS.RegisterCodeExpire)
	}
	if cfg.SMS.DevReturnCode {
		t.Fatal("default config should not return local register code")
	}
	if cfg.RPC.Trip.ServiceName != "trip-service" {
		t.Fatalf("unexpected trip service name: %s", cfg.RPC.Trip.ServiceName)
	}
	if cfg.RPC.Trip.Port != 9003 {
		t.Fatalf("unexpected trip service port: %d", cfg.RPC.Trip.Port)
	}
	if cfg.RPC.Trip.Target != "127.0.0.1:9003" {
		t.Fatalf("unexpected trip service target: %s", cfg.RPC.Trip.Target)
	}
	if !cfg.TravelData.Enabled {
		t.Fatal("travel data should be enabled by default")
	}
	if cfg.TravelData.Timeout != 5*time.Second {
		t.Fatalf("unexpected travel data timeout: %s", cfg.TravelData.Timeout)
	}
}

func TestLoadSMSDevReturnCodeFromConfigFile(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	content := []byte("sms:\n  dev_return_code: true\n")
	if err := os.WriteFile(configPath, content, 0600); err != nil {
		t.Fatalf("write temp config failed: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}
	if !cfg.SMS.DevReturnCode {
		t.Fatal("sms dev_return_code should be loaded from config file")
	}
}

func TestLoadTravelDataConfigFromConfigFile(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	content := []byte("travel_data:\n  enabled: true\n  amap_key: amap-test-key\n  qweather_key: qweather-test-key\n  timeout: 8s\n")
	if err := os.WriteFile(configPath, content, 0600); err != nil {
		t.Fatalf("write temp config failed: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}
	if !cfg.TravelData.Enabled {
		t.Fatal("travel data should be enabled from config file")
	}
	if cfg.TravelData.AmapKey != "amap-test-key" {
		t.Fatalf("unexpected amap key: %s", cfg.TravelData.AmapKey)
	}
	if cfg.TravelData.QWeatherKey != "qweather-test-key" {
		t.Fatalf("unexpected qweather key: %s", cfg.TravelData.QWeatherKey)
	}
	if cfg.TravelData.Timeout != 8*time.Second {
		t.Fatalf("unexpected travel data timeout: %s", cfg.TravelData.Timeout)
	}
}
