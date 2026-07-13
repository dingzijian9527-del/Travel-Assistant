package config

import (
	"testing"
)

func TestMaskString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "", ""},
		{"single char", "a", "a"},
		{"two chars", "ab", "ab"},
		{"three chars", "abc", "abc"},
		{"four chars", "abcd", "abcd"},
		{"five chars", "abcde", "ab***de"},
		{"six chars", "abcdef", "ab***ef"},
		{"long string", "hello world", "he***ld"},
		{"chinese five chars", "一二三四五", "一二***四五"},
		{"chinese six chars", "一二三四五六", "一二***五六"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskString(tt.input)
			if result != tt.expected {
				t.Errorf("MaskString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMaskDSN(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard DSN",
			input:    "root:password@tcp(127.0.0.1:3306)/travel-assistant?charset=utf8mb4&parseTime=True&loc=Local",
			expected: "root:***@tcp(127.0.0.1:3306)/travel-assistant?charset=utf8mb4&parseTime=True&loc=Local",
		},
		{
			name:     "no password - only colon",
			input:    "root:@tcp(127.0.0.1:3306)/db",
			expected: "root:***@tcp(127.0.0.1:3306)/db",
		},
		{
			name:     "no at sign",
			input:    "just a string",
			expected: "just a string",
		},
		{
			name:     "complex password with special chars",
			input:    "admin:P@ssw0rd!@tcp(localhost:3306)/mydb",
			expected: "admin:***@tcp(localhost:3306)/mydb",
		},
		{
			name:     "at sign in params",
			input:    "user:pass@tcp(host:3306)/db?foo=@bar",
			expected: "user:***@tcp(host:3306)/db?foo=@bar",
		},
		{
			name:     "empty DSN",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskDSN(tt.input)
			if result != tt.expected {
				t.Errorf("MaskDSN(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSafeConfigSummary(t *testing.T) {
	cfg := &Config{
		App: AppConfig{
			Env: "production",
		},
		HTTP: HTTPConfig{
			Port: 8080,
		},
		MySQL: MySQLConfig{
			DSN: "root:secret@tcp(127.0.0.1:3306)/db",
		},
		Redis: RedisConfig{
			Addr:     "127.0.0.1:6379",
			Password: "redis-secret",
		},
		RAG: RAGConfig{
			Enabled: true,
		},
		TravelData: TravelDataConfig{
			Enabled: false,
		},
	}

	result := SafeConfigSummary(cfg)

	// Check expected keys and values
	if result["env"] != "production" {
		t.Errorf("env = %v, want production", result["env"])
	}
	if result["http_port"] != 8080 {
		t.Errorf("http_port = %v, want 8080", result["http_port"])
	}
	if result["mysql_enabled"] != true {
		t.Errorf("mysql_enabled = %v, want true", result["mysql_enabled"])
	}
	if result["redis_enabled"] != true {
		t.Errorf("redis_enabled = %v, want true", result["redis_enabled"])
	}
	if result["rag_enabled"] != true {
		t.Errorf("rag_enabled = %v, want true", result["rag_enabled"])
	}
	if result["travel_data_enabled"] != false {
		t.Errorf("travel_data_enabled = %v, want false", result["travel_data_enabled"])
	}

	// Verify sensitive fields are NOT in the summary
	sensitiveKeys := []string{"api_key", "jwt_secret", "secret_key", "password", "dsn", "amap_key", "qweather_key"}
	for _, key := range sensitiveKeys {
		if _, exists := result[key]; exists {
			t.Errorf("sensitive key %q should not be in SafeConfigSummary", key)
		}
	}
}

func TestSafeConfigSummaryDisabledServices(t *testing.T) {
	cfg := &Config{
		App:  AppConfig{Env: "dev"},
		HTTP: HTTPConfig{Port: 3000},
	}

	result := SafeConfigSummary(cfg)

	if result["mysql_enabled"] != false {
		t.Errorf("mysql_enabled should be false when DSN is empty")
	}
	if result["redis_enabled"] != false {
		t.Errorf("redis_enabled should be false when Addr is empty")
	}
	if result["rag_enabled"] != false {
		t.Errorf("rag_enabled should be false when RAG.Enabled is false")
	}
	if result["travel_data_enabled"] != false {
		t.Errorf("travel_data_enabled should be false when TravelData.Enabled is false")
	}
}

func TestMaskStringInSafeConfigSummary(t *testing.T) {
	// Verify that SafeConfigSummary never leaks sensitive values
	cfg := &Config{
		App:  AppConfig{Env: "staging"},
		HTTP: HTTPConfig{Port: 9090},
		MySQL: MySQLConfig{
			DSN: "admin:super-secret-password@tcp(db.example.com:3306)/production",
		},
	}

	result := SafeConfigSummary(cfg)

	// The summary should NOT contain any raw password or DSN
	for k, v := range result {
		strVal, ok := v.(string)
		if !ok {
			continue
		}
		if strVal == "super-secret-password" {
			t.Errorf("key %q leaks password: %q", k, strVal)
		}
	}

	if result["env"] != "staging" {
		t.Errorf("env = %v, want staging", result["env"])
	}
	if result["mysql_enabled"] != true {
		t.Errorf("mysql_enabled should be true")
	}
}
