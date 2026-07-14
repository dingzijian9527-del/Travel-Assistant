package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadDefaultRAGConfig(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(configPath, []byte(""), 0600); err != nil {
		t.Fatalf("write temp config failed: %v", err)
	}

	cfg, err := Load(configPath)
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
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(configPath, []byte(""), 0600); err != nil {
		t.Fatalf("write temp config failed: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}

	if cfg.MySQL.DSN != "" {
		t.Fatalf("mysql dsn default should be empty, got: %s", cfg.MySQL.DSN)
	}
	if cfg.Redis.Addr == "" {
		t.Fatal("redis addr should be configured")
	}
	if cfg.Auth.JWTSecret != "" {
		t.Fatalf("jwt secret default should be empty, got: %q", cfg.Auth.JWTSecret)
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

func TestLoadUsesLocalConfigBeforeExampleConfig(t *testing.T) {
	root := t.TempDir()
	confDir := filepath.Join(root, "conf")
	if err := os.MkdirAll(confDir, 0o755); err != nil {
		t.Fatalf("create conf dir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(confDir, "config.yaml"), []byte("auth:\n  jwt_secret: example-secret\n"), 0o600); err != nil {
		t.Fatalf("write example config failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(confDir, "config.local.yaml"), []byte("auth:\n  jwt_secret: local-secret\n"), 0o600); err != nil {
		t.Fatalf("write local config failed: %v", err)
	}
	t.Chdir(root)

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}
	if cfg.Auth.JWTSecret != "local-secret" {
		t.Fatalf("jwt secret = %q, want local-secret", cfg.Auth.JWTSecret)
	}
}

func TestLoadExplicitPathHasHighestPriority(t *testing.T) {
	root := t.TempDir()
	confDir := filepath.Join(root, "conf")
	if err := os.MkdirAll(confDir, 0o755); err != nil {
		t.Fatalf("create conf dir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(confDir, "config.local.yaml"), []byte("auth:\n  jwt_secret: local-secret\n"), 0o600); err != nil {
		t.Fatalf("write local config failed: %v", err)
	}
	explicitPath := filepath.Join(root, "explicit.yaml")
	if err := os.WriteFile(explicitPath, []byte("auth:\n  jwt_secret: explicit-secret\n"), 0o600); err != nil {
		t.Fatalf("write explicit config failed: %v", err)
	}
	t.Chdir(root)

	cfg, err := Load(explicitPath)
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}
	if cfg.Auth.JWTSecret != "explicit-secret" {
		t.Fatalf("jwt secret = %q, want explicit-secret", cfg.Auth.JWTSecret)
	}
}

func TestLoadReturnsErrorForInvalidLocalConfig(t *testing.T) {
	root := t.TempDir()
	confDir := filepath.Join(root, "conf")
	if err := os.MkdirAll(confDir, 0o755); err != nil {
		t.Fatalf("create conf dir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(confDir, "config.yaml"), []byte("auth:\n  jwt_secret: example-secret\n"), 0o600); err != nil {
		t.Fatalf("write example config failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(confDir, "config.local.yaml"), []byte("auth: [\n"), 0o600); err != nil {
		t.Fatalf("write invalid local config failed: %v", err)
	}
	t.Chdir(root)

	if _, err := Load(""); err == nil {
		t.Fatal("invalid local config should return an error")
	}
}

func TestLoadDoesNotOverrideFileValueFromEnvironment(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(configPath, []byte("auth:\n  jwt_secret: file-secret\n"), 0o600); err != nil {
		t.Fatalf("write config failed: %v", err)
	}
	t.Setenv("TRAVEL_ASSISTANT_AUTH_JWT_SECRET", "environment-secret")

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}
	if cfg.Auth.JWTSecret != "file-secret" {
		t.Fatalf("jwt secret = %q, want file-secret", cfg.Auth.JWTSecret)
	}
}
func TestValidateForServiceAllowsDevelopmentConfigWithoutOptionalThirdPartyKeys(t *testing.T) {
	cfg := validConfigForService("api-gateway")

	if err := cfg.ValidateForService("api-gateway"); err != nil {
		t.Fatalf("development config should be valid: %v", err)
	}
}

func TestValidateForServiceRejectsPlaceholderJWTSecret(t *testing.T) {
	cfg := validConfigForService("api-gateway")
	cfg.Auth.JWTSecret = "change-me-in-local-config"

	err := cfg.ValidateForService("api-gateway")
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "jwt") {
		t.Fatalf("expected jwt validation error, got: %v", err)
	}
}

func TestValidateForServiceRejectsProductionUnsafeSMSFallback(t *testing.T) {
	cfg := validConfigForService("api-gateway")
	cfg.App.Env = "prod"
	cfg.SMS.DevReturnCode = true

	err := cfg.ValidateForService("api-gateway")
	if err == nil || !strings.Contains(err.Error(), "dev_return_code") {
		t.Fatalf("expected sms dev_return_code validation error, got: %v", err)
	}
}

func TestValidateForServiceRejectsProductionWithoutAIKey(t *testing.T) {
	cfg := validConfigForService("api-gateway")
	cfg.App.Env = "prod"
	cfg.AI.APIKey = ""
	cfg.AI.EndpointID = ""
	cfg.Upload.Qiniu = UploadQiniuConfig{}
	cfg.TravelData.AmapKey = ""
	cfg.TravelData.QWeatherKey = ""

	err := cfg.ValidateForService("api-gateway")
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "ai") {
		t.Fatalf("expected ai validation error, got: %v", err)
	}
}

func validConfigForService(service string) *Config {
	cfg := &Config{
		App:  AppConfig{Name: "travel-assistant", Env: "dev"},
		HTTP: HTTPConfig{Host: "127.0.0.1", Port: 8080, ReadTimeout: 5 * time.Second, WriteTimeout: 10 * time.Second},
		RPC: RPCConfig{
			Host:    "127.0.0.1",
			User:    RPCServiceConfig{ServiceName: "user-service", Port: 9001, Target: "127.0.0.1:9001"},
			AIAgent: RPCServiceConfig{ServiceName: "ai-agent-service", Port: 9002, Target: "127.0.0.1:9002"},
			Trip:    RPCServiceConfig{ServiceName: "trip-service", Port: 9003, Target: "127.0.0.1:9003"},
		},
		MySQL: MySQLConfig{DSN: "tester:secret@tcp(127.0.0.1:3306)/travel-assistant"},
		Redis: RedisConfig{Addr: "127.0.0.1:6379"},
		SMS:   SMSConfig{Endpoint: "https://sms.tencentcloudapi.com", RegisterCodeExpire: 5 * time.Minute},
		Auth:  AuthConfig{JWTSecret: "unit-test-jwt-secret", JWTExpire: 24 * time.Hour},
		Upload: UploadConfig{
			LocalDir:          "uploads",
			MaxSizeMB:         20,
			AllowedExtensions: []string{".jpg"},
		},
		AI:         AIConfig{Provider: "ark", MaxPromptChars: 2000},
		RAG:        RAGConfig{Enabled: true, Address: "127.0.0.1:19530", TopK: 3, MinScore: 0.15, EmbeddingDim: 768, CollectionName: "travel_knowledge", Provider: "local"},
		TravelData: TravelDataConfig{Enabled: true, Timeout: 5 * time.Second},
	}

	switch service {
	case "user-service", "trip-service", "ai-agent-service":
		cfg.HTTP = HTTPConfig{}
	}
	return cfg
}
