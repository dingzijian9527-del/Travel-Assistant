package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitWithConfigPathRejectsUnsafeJWTPlaceholder(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	logDir, err := os.MkdirTemp("", "travel-assistant-bootstrap-log-*")
	if err != nil {
		t.Fatalf("create log dir failed: %v", err)
	}
	content := []byte(fmt.Sprintf("app:\n  name: travel-assistant\n  env: dev\nhttp:\n  host: 127.0.0.1\n  port: 8080\nlog:\n  dir: %q\nmysql:\n  dsn: tester:secret@tcp(127.0.0.1:3306)/travel-assistant\nredis:\n  addr: 127.0.0.1:6379\nauth:\n  jwt_secret: change-me-in-local-config\n", strings.ReplaceAll(logDir, "\\", "/")))
	if err := os.WriteFile(configPath, content, 0600); err != nil {
		t.Fatalf("write config failed: %v", err)
	}

	_, err = InitWithConfigPath("api-gateway", configPath)
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "jwt") {
		t.Fatalf("expected jwt validation error, got: %v", err)
	}
}

func TestInitWithConfigPathLoadsRuntimeWhenConfigIsValid(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	logDir, err := os.MkdirTemp("", "travel-assistant-bootstrap-log-*")
	if err != nil {
		t.Fatalf("create log dir failed: %v", err)
	}
	content := []byte(fmt.Sprintf("app:\n  name: travel-assistant\n  env: dev\nhttp:\n  host: 127.0.0.1\n  port: 8080\nlog:\n  dir: %q\nmysql:\n  dsn: tester:secret@tcp(127.0.0.1:3306)/travel-assistant\nredis:\n  addr: 127.0.0.1:6379\nauth:\n  jwt_secret: unit-test-jwt-secret\n", strings.ReplaceAll(logDir, "\\", "/")))
	if err := os.WriteFile(configPath, content, 0600); err != nil {
		t.Fatalf("write config failed: %v", err)
	}

	runtime, err := InitWithConfigPath("api-gateway", configPath)
	if err != nil {
		t.Fatalf("expected init success, got: %v", err)
	}
	if runtime == nil || runtime.Config == nil || runtime.Logger == nil {
		t.Fatalf("expected runtime with config and logger, got: %#v", runtime)
	}
	_ = runtime.Logger.Sync()
}
