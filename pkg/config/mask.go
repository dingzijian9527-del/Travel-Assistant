package config

import (
	"fmt"
	"strings"
)

// MaskString masks the middle portion of a string, keeping the first 2 and last 2 characters.
// Strings shorter than 5 characters (in runes) are returned unchanged.
func MaskString(s string) string {
	runes := []rune(s)
	if len(runes) < 5 {
		return s
	}
	return string(runes[:2]) + "***" + string(runes[len(runes)-2:])
}

// MaskDSN masks the password portion of a MySQL DSN string.
// Typical DSN format: user:password@tcp(host:port)/dbname?params
func MaskDSN(dsn string) string {
	// Locate the credential separator: the '@' that immediately precedes a protocol like 'tcp('.
	// We search for '@tcp(' as the canonical MySQL driver DSN separator.
	atIdx := strings.Index(dsn, "@tcp(")
	if atIdx == -1 {
		// Fallback: search for '@' that has a protocol+paren after it
		atIdx = strings.Index(dsn, "@")
		if atIdx == -1 {
			return dsn
		}
	}
	colonIdx := strings.Index(dsn[:atIdx], ":")
	if colonIdx == -1 {
		return dsn
	}
	return dsn[:colonIdx+1] + "***" + dsn[atIdx:]
}

// SafeConfigSummary returns a sanitized configuration summary without any sensitive values.
// Only exposes: env, http_port, mysql_enabled, redis_enabled, rag_enabled, travel_data_enabled.
// api_key, jwt_secret, secret_key, and other sensitive fields are never exposed.
func SafeConfigSummary(cfg *Config) map[string]any {
	mysqlEnabled := cfg.MySQL.DSN != ""
	redisEnabled := cfg.Redis.Addr != ""

	return map[string]any{
		"env":                 cfg.App.Env,
		"http_port":           cfg.HTTP.Port,
		"mysql_enabled":       mysqlEnabled,
		"redis_enabled":       redisEnabled,
		"rag_enabled":         cfg.RAG.Enabled,
		"travel_data_enabled": cfg.TravelData.Enabled,
	}
}

// String returns a human-readable masked representation of the Config (for logging).
func (c *Config) String() string {
	return fmt.Sprintf(
		"Config{env=%s, http=%s:%d, mysql=%s, redis=%s, rag=%s, travel_data=%s}",
		c.App.Env,
		c.HTTP.Host,
		c.HTTP.Port,
		MaskDSN(c.MySQL.DSN),
		c.Redis.Addr,
		func() string {
			if c.RAG.Enabled {
				return "enabled"
			}
			return "disabled"
		}(),
		func() string {
			if c.TravelData.Enabled {
				return "enabled"
			}
			return "disabled"
		}(),
	)
}
