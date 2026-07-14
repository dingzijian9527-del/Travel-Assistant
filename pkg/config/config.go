package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type HTTPConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type RPCServiceConfig struct {
	ServiceName string `mapstructure:"service_name"`
	Port        int    `mapstructure:"port"`
	Target      string `mapstructure:"target"`
}

type RPCConfig struct {
	Host    string           `mapstructure:"host"`
	User    RPCServiceConfig `mapstructure:"user"`
	AIAgent RPCServiceConfig `mapstructure:"ai_agent"`
	Trip    RPCServiceConfig `mapstructure:"trip"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Dir        string `mapstructure:"dir"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
	Console    bool   `mapstructure:"console"`
}

type MySQLConfig struct {
	DSN                    string `mapstructure:"dsn"`
	MaxIdleConns           int    `mapstructure:"max_idle_conns"`
	MaxOpenConns           int    `mapstructure:"max_open_conns"`
	ConnMaxLifetimeSeconds int    `mapstructure:"conn_max_lifetime_seconds"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type SMSConfig struct {
	SecretID           string        `mapstructure:"secret_id"`
	SecretKey          string        `mapstructure:"secret_key"`
	SDKAppID           string        `mapstructure:"sdk_app_id"`
	SignName           string        `mapstructure:"sign_name"`
	TemplateID         string        `mapstructure:"template_id"`
	Region             string        `mapstructure:"region"`
	Endpoint           string        `mapstructure:"endpoint"`
	RegisterCodeExpire time.Duration `mapstructure:"register_code_expire"`
	DevReturnCode      bool          `mapstructure:"dev_return_code"`
}

type AuthConfig struct {
	JWTSecret        string        `mapstructure:"jwt_secret"`
	JWTExpire        time.Duration `mapstructure:"jwt_expire"`
	JWTRefreshExpire time.Duration `mapstructure:"jwt_refresh_expire"`
}

type UploadQiniuConfig struct {
	AccessKey  string `mapstructure:"access_key"`
	SecretKey  string `mapstructure:"secret_key"`
	Bucket     string `mapstructure:"bucket"`
	URL        string `mapstructure:"url"`
	UploadHost string `mapstructure:"upload_host"`
}

type UploadConfig struct {
	LocalDir          string            `mapstructure:"local_dir"`
	MaxSizeMB         int64             `mapstructure:"max_size_mb"`
	AllowedExtensions []string          `mapstructure:"allowed_extensions"`
	Qiniu             UploadQiniuConfig `mapstructure:"qiniu"`
}

type AIConfig struct {
	Provider       string        `mapstructure:"provider"`
	APIKey         string        `mapstructure:"api_key"`
	BaseURL        string        `mapstructure:"base_url"`
	EndpointID     string        `mapstructure:"endpoint_id"`
	ModelName      string        `mapstructure:"model_name"`
	Model          string        `mapstructure:"model"`
	Timeout        time.Duration `mapstructure:"timeout"`
	Stream         bool          `mapstructure:"stream"`
	SystemPrompt   string        `mapstructure:"system_prompt"`
	MaxPromptChars int           `mapstructure:"max_prompt_chars"`
}

type RAGConfig struct {
	Enabled          bool    `mapstructure:"enabled"`
	Provider         string  `mapstructure:"provider"`
	Address          string  `mapstructure:"address"`
	CollectionName   string  `mapstructure:"collection_name"`
	EmbeddingDim     int     `mapstructure:"embedding_dim"`
	EmbeddingBaseURL string  `mapstructure:"embedding_base_url"`
	EmbeddingAPIKey  string  `mapstructure:"embedding_api_key"`
	EmbeddingModel   string  `mapstructure:"embedding_model"`
	TopK             int     `mapstructure:"top_k"`
	MinScore         float64 `mapstructure:"min_score"`
}

type TravelDataConfig struct {
	Enabled     bool          `mapstructure:"enabled"`
	AmapKey     string        `mapstructure:"amap_key"`
	QWeatherKey string        `mapstructure:"qweather_key"`
	Timeout     time.Duration `mapstructure:"timeout"`
}

type Config struct {
	App        AppConfig        `mapstructure:"app"`
	HTTP       HTTPConfig       `mapstructure:"http"`
	RPC        RPCConfig        `mapstructure:"rpc"`
	Log        LogConfig        `mapstructure:"log"`
	MySQL      MySQLConfig      `mapstructure:"mysql"`
	Redis      RedisConfig      `mapstructure:"redis"`
	SMS        SMSConfig        `mapstructure:"sms"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Upload     UploadConfig     `mapstructure:"upload"`
	AI         AIConfig         `mapstructure:"ai"`
	RAG        RAGConfig        `mapstructure:"rag"`
	TravelData TravelDataConfig `mapstructure:"travel_data"`
}

const placeholderJWTSecret = "change-me-in-local-config"

var global *Config

func Load(path string) (*Config, error) {
	v := viper.New()
	setDefaults(v)
	v.SetConfigType("yaml")
	if strings.TrimSpace(path) != "" {
		v.SetConfigFile(path)
	} else {
		v.SetConfigName("config")
		v.AddConfigPath("conf")
		v.AddConfigPath(".")
	}
	v.SetEnvPrefix("TRAVEL_ASSISTANT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func InitGlobal(path string) (*Config, error) {
	cfg, err := Load(path)
	if err != nil {
		return nil, err
	}
	global = cfg
	return cfg, nil
}

func MustGlobal() *Config {
	if global == nil {
		panic("全局配置尚未初始化")
	}
	return global
}

func (c *Config) ValidateForService(service string) error {
	if c == nil {
		return fmt.Errorf("config is required")
	}
	service = strings.TrimSpace(service)
	if strings.TrimSpace(c.App.Name) == "" {
		return fmt.Errorf("app.name is required")
	}
	if service == "api-gateway" {
		if strings.TrimSpace(c.HTTP.Host) == "" || c.HTTP.Port <= 0 {
			return fmt.Errorf("http.host and http.port are required for %s", service)
		}
	}
	if err := validateJWTSecret(c.Auth.JWTSecret); err != nil {
		return err
	}
	if requiresMySQL(service) && strings.TrimSpace(c.MySQL.DSN) == "" {
		return fmt.Errorf("mysql.dsn is required for %s", service)
	}
	if requiresRPCHost(service) && strings.TrimSpace(c.RPC.Host) == "" {
		return fmt.Errorf("rpc.host is required for %s", service)
	}

	switch service {
	case "user-service":
		if err := validateRPCServiceConfig("rpc.user", c.RPC.User); err != nil {
			return err
		}
	case "ai-agent-service":
		if err := validateRPCServiceConfig("rpc.ai_agent", c.RPC.AIAgent); err != nil {
			return err
		}
	case "trip-service":
		if err := validateRPCServiceConfig("rpc.trip", c.RPC.Trip); err != nil {
			return err
		}
	case "api-gateway", "":
		if strings.TrimSpace(c.Redis.Addr) == "" {
			return fmt.Errorf("redis.addr is required for api-gateway")
		}
	}

	if !isProductionEnv(c.App.Env) {
		return nil
	}
	if c.SMS.DevReturnCode {
		return fmt.Errorf("sms.dev_return_code must be false in production")
	}
	if service == "api-gateway" || service == "ai-agent-service" {
		if strings.TrimSpace(c.AI.APIKey) == "" {
			return fmt.Errorf("ai.api_key is required in production")
		}
		if strings.TrimSpace(c.AI.EndpointID) == "" && strings.TrimSpace(c.AI.ModelName) == "" && strings.TrimSpace(c.AI.Model) == "" {
			return fmt.Errorf("ai endpoint or model is required in production")
		}
	}
	if service == "api-gateway" {
		if !tencentSMSConfigComplete(c.SMS) {
			return fmt.Errorf("sms config is required in production")
		}
		if !qiniuConfigComplete(c.Upload.Qiniu) {
			return fmt.Errorf("upload.qiniu config is required in production")
		}
		if c.TravelData.Enabled && (strings.TrimSpace(c.TravelData.AmapKey) == "" || strings.TrimSpace(c.TravelData.QWeatherKey) == "") {
			return fmt.Errorf("travel_data keys are required in production")
		}
	}
	if c.RAG.Enabled {
		provider := strings.ToLower(strings.TrimSpace(c.RAG.Provider))
		if provider == "" || provider == "local" {
			return fmt.Errorf("rag.provider must use a semantic embedder in production")
		}
		if strings.TrimSpace(c.RAG.EmbeddingBaseURL) == "" || strings.TrimSpace(c.RAG.EmbeddingAPIKey) == "" || strings.TrimSpace(c.RAG.EmbeddingModel) == "" {
			return fmt.Errorf("semantic embedding config is required in production")
		}
	}
	return nil
}

func normalizeEnv(env string) string {
	env = strings.ToLower(strings.TrimSpace(env))
	if env == "" {
		return "dev"
	}
	return env
}

func isProductionEnv(env string) bool {
	env = normalizeEnv(env)
	return env == "prod" || env == "production"
}

func requiresMySQL(service string) bool {
	switch service {
	case "api-gateway", "user-service", "ai-agent-service", "trip-service", "":
		return true
	default:
		return false
	}
}

func requiresRPCHost(service string) bool {
	switch service {
	case "user-service", "ai-agent-service", "trip-service":
		return true
	default:
		return false
	}
}

func validateJWTSecret(secret string) error {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return fmt.Errorf("auth.jwt_secret is required")
	}
	if secret == placeholderJWTSecret {
		return fmt.Errorf("auth.jwt_secret must not use placeholder value")
	}
	return nil
}

func validateRPCServiceConfig(field string, cfg RPCServiceConfig) error {
	if strings.TrimSpace(cfg.ServiceName) == "" {
		return fmt.Errorf("%s.service_name is required", field)
	}
	if cfg.Port <= 0 {
		return fmt.Errorf("%s.port must be positive", field)
	}
	return nil
}

func tencentSMSConfigComplete(cfg SMSConfig) bool {
	return strings.TrimSpace(cfg.SecretID) != "" &&
		strings.TrimSpace(cfg.SecretKey) != "" &&
		strings.TrimSpace(cfg.SDKAppID) != "" &&
		strings.TrimSpace(cfg.SignName) != "" &&
		strings.TrimSpace(cfg.TemplateID) != ""
}

func qiniuConfigComplete(cfg UploadQiniuConfig) bool {
	return strings.TrimSpace(cfg.AccessKey) != "" &&
		strings.TrimSpace(cfg.SecretKey) != "" &&
		strings.TrimSpace(cfg.Bucket) != "" &&
		strings.TrimSpace(cfg.URL) != ""
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "travel-assistant")
	v.SetDefault("app.env", "dev")
	v.SetDefault("http.host", "0.0.0.0")
	v.SetDefault("http.port", 8080)
	v.SetDefault("http.read_timeout", "5s")
	v.SetDefault("http.write_timeout", "10s")
	v.SetDefault("rpc.host", "0.0.0.0")
	v.SetDefault("rpc.user.service_name", "user-service")
	v.SetDefault("rpc.user.port", 9001)
	v.SetDefault("rpc.user.target", "127.0.0.1:9001")
	v.SetDefault("rpc.ai_agent.service_name", "ai-agent-service")
	v.SetDefault("rpc.ai_agent.port", 9002)
	v.SetDefault("rpc.ai_agent.target", "127.0.0.1:9002")
	v.SetDefault("rpc.trip.service_name", "trip-service")
	v.SetDefault("rpc.trip.port", 9003)
	v.SetDefault("rpc.trip.target", "127.0.0.1:9003")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.dir", "logs")
	v.SetDefault("log.max_age_days", 7)
	v.SetDefault("log.console", true)
	v.SetDefault("mysql.dsn", "")
	v.SetDefault("mysql.max_idle_conns", 10)
	v.SetDefault("mysql.max_open_conns", 50)
	v.SetDefault("mysql.conn_max_lifetime_seconds", 3600)
	v.SetDefault("redis.addr", "127.0.0.1:6379")
	v.SetDefault("redis.username", "")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("sms.secret_id", "")
	v.SetDefault("sms.secret_key", "")
	v.SetDefault("sms.sdk_app_id", "")
	v.SetDefault("sms.sign_name", "")
	v.SetDefault("sms.template_id", "")
	v.SetDefault("sms.region", "ap-guangzhou")
	v.SetDefault("sms.endpoint", "https://sms.tencentcloudapi.com")
	v.SetDefault("sms.register_code_expire", "5m")
	v.SetDefault("sms.dev_return_code", false)
	v.SetDefault("auth.jwt_secret", "")
	v.SetDefault("auth.jwt_expire", "24h")
	v.SetDefault("auth.jwt_refresh_expire", "168h")
	v.SetDefault("upload.local_dir", "uploads")
	v.SetDefault("upload.max_size_mb", 20)
	v.SetDefault("upload.allowed_extensions", []string{".jpg", ".jpeg", ".png", ".webp", ".pdf"})
	v.SetDefault("upload.qiniu.access_key", "")
	v.SetDefault("upload.qiniu.secret_key", "")
	v.SetDefault("upload.qiniu.bucket", "")
	v.SetDefault("upload.qiniu.url", "")
	v.SetDefault("upload.qiniu.upload_host", "https://up-z2.qiniup.com")
	v.SetDefault("ai.provider", "ark")
	v.SetDefault("ai.api_key", "")
	v.SetDefault("ai.base_url", "https://ark.cn-beijing.volces.com/api/v3")
	v.SetDefault("ai.endpoint_id", "")
	v.SetDefault("ai.model_name", "")
	v.SetDefault("ai.model", "")
	v.SetDefault("ai.timeout", "60s")
	v.SetDefault("ai.stream", true)
	v.SetDefault("ai.system_prompt", "你是旅行助手项目中的专属旅行规划智能体，只服务旅游出行场景。")
	v.SetDefault("ai.max_prompt_chars", 2000)
	v.SetDefault("rag.enabled", true)
	v.SetDefault("rag.provider", "local")
	v.SetDefault("rag.address", "127.0.0.1:19530")
	v.SetDefault("rag.collection_name", "travel_knowledge")
	v.SetDefault("rag.embedding_dim", 768)
	v.SetDefault("rag.embedding_base_url", "")
	v.SetDefault("rag.embedding_api_key", "")
	v.SetDefault("rag.embedding_model", "")
	v.SetDefault("rag.top_k", 3)
	v.SetDefault("rag.min_score", 0.15)
	v.SetDefault("travel_data.enabled", true)
	v.SetDefault("travel_data.amap_key", "")
	v.SetDefault("travel_data.qweather_key", "")
	v.SetDefault("travel_data.timeout", "5s")
}
