package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// AppConfig 描述应用基础信息。
type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

// HTTPConfig 描述网关监听与超时配置。
type HTTPConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// RPCServiceConfig 描述单个远程过程调用服务的本地端口和调用地址。
type RPCServiceConfig struct {
	ServiceName string `mapstructure:"service_name"`
	Port        int    `mapstructure:"port"`
	Target      string `mapstructure:"target"`
}

// RPCConfig 描述所有远程过程调用服务的统一配置入口。
type RPCConfig struct {
	Host    string           `mapstructure:"host"`
	User    RPCServiceConfig `mapstructure:"user"`
	AIAgent RPCServiceConfig `mapstructure:"ai_agent"`
}

// LogConfig 描述日志输出和保留策略。
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Dir        string `mapstructure:"dir"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
	Console    bool   `mapstructure:"console"`
}

// MySQLConfig 描述用户服务持久化所需的 MySQL 配置。
type MySQLConfig struct {
	DSN                    string `mapstructure:"dsn"`
	MaxIdleConns           int    `mapstructure:"max_idle_conns"`
	MaxOpenConns           int    `mapstructure:"max_open_conns"`
	ConnMaxLifetimeSeconds int    `mapstructure:"conn_max_lifetime_seconds"`
}

// RedisConfig 描述缓存服务连接配置，当前预留给会话和热点数据缓存。
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// SMSConfig 描述腾讯云短信验证码发送配置。
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

// AuthConfig 描述网关和用户服务共享的认证令牌配置。
type AuthConfig struct {
	JWTSecret string        `mapstructure:"jwt_secret"`
	JWTExpire time.Duration `mapstructure:"jwt_expire"`
}

// UploadConfig 描述文件上传限制。
type UploadConfig struct {
	LocalDir          string   `mapstructure:"local_dir"`
	MaxSizeMB         int64    `mapstructure:"max_size_mb"`
	AllowedExtensions []string `mapstructure:"allowed_extensions"`
}

// AIConfig 描述旅行智能体调用大模型所需配置。
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

// RAGConfig describes retrieval augmented generation settings.
type RAGConfig struct {
	Enabled        bool    `mapstructure:"enabled"`
	Provider       string  `mapstructure:"provider"`
	Address        string  `mapstructure:"address"`
	CollectionName string  `mapstructure:"collection_name"`
	EmbeddingDim   int     `mapstructure:"embedding_dim"`
	TopK           int     `mapstructure:"top_k"`
	MinScore       float64 `mapstructure:"min_score"`
}

// Config 是进程启动时加载的完整配置。
type Config struct {
	App    AppConfig    `mapstructure:"app"`
	HTTP   HTTPConfig   `mapstructure:"http"`
	RPC    RPCConfig    `mapstructure:"rpc"`
	Log    LogConfig    `mapstructure:"log"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Redis  RedisConfig  `mapstructure:"redis"`
	SMS    SMSConfig    `mapstructure:"sms"`
	Auth   AuthConfig   `mapstructure:"auth"`
	Upload UploadConfig `mapstructure:"upload"`
	AI     AIConfig     `mapstructure:"ai"`
	RAG    RAGConfig    `mapstructure:"rag"`
}

var global *Config

// Load 使用配置读取工具从配置文件和环境变量加载配置。
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

// InitGlobal 加载并缓存进程级配置。
func InitGlobal(path string) (*Config, error) {
	cfg, err := Load(path)
	if err != nil {
		return nil, err
	}
	global = cfg
	return cfg, nil
}

// MustGlobal 返回进程级配置，未初始化时直接失败。
func MustGlobal() *Config {
	if global == nil {
		panic("配置未初始化")
	}
	return global
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
	v.SetDefault("log.level", "info")
	v.SetDefault("log.dir", "logs")
	v.SetDefault("log.max_age_days", 7)
	v.SetDefault("log.console", true)
	v.SetDefault("mysql.dsn", "root:123456@tcp(127.0.0.1:3306)/travel-assistant?charset=utf8mb4&parseTime=True&loc=Local")
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
	v.SetDefault("auth.jwt_secret", "travel-assistant-dev-secret")
	v.SetDefault("auth.jwt_expire", "24h")
	v.SetDefault("upload.local_dir", "uploads")
	v.SetDefault("upload.max_size_mb", 20)
	v.SetDefault("upload.allowed_extensions", []string{".jpg", ".jpeg", ".png", ".webp", ".pdf"})
	v.SetDefault("ai.provider", "ark")
	v.SetDefault("ai.api_key", "")
	v.SetDefault("ai.base_url", "https://ark.cn-beijing.volces.com/api/v3")
	v.SetDefault("ai.endpoint_id", "")
	v.SetDefault("ai.model_name", "")
	v.SetDefault("ai.model", "")
	v.SetDefault("ai.timeout", "60s")
	v.SetDefault("ai.stream", true)
	v.SetDefault("ai.system_prompt", "你是旅行助手项目中的旅游专用智能体，只回答出行计划、旅游攻略、目的地、交通、酒店、天气和美食相关问题。")
	v.SetDefault("ai.max_prompt_chars", 2000)
	v.SetDefault("rag.enabled", true)
	v.SetDefault("rag.provider", "local")
	v.SetDefault("rag.address", "115.190.209.83:19530")
	v.SetDefault("rag.collection_name", "travel_knowledge")
	v.SetDefault("rag.embedding_dim", 768)
	v.SetDefault("rag.top_k", 3)
	v.SetDefault("rag.min_score", 0.15)
}
