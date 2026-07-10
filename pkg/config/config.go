package config

import (
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
	JWTSecret string        `mapstructure:"jwt_secret"`
	JWTExpire time.Duration `mapstructure:"jwt_expire"`
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
	Enabled        bool    `mapstructure:"enabled"`
	Provider       string  `mapstructure:"provider"`
	Address        string  `mapstructure:"address"`
	CollectionName string  `mapstructure:"collection_name"`
	EmbeddingDim   int     `mapstructure:"embedding_dim"`
	TopK           int     `mapstructure:"top_k"`
	MinScore       float64 `mapstructure:"min_score"`
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
		panic("閰嶇疆鏈垵濮嬪寲")
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
	v.SetDefault("rpc.trip.service_name", "trip-service")
	v.SetDefault("rpc.trip.port", 9003)
	v.SetDefault("rpc.trip.target", "127.0.0.1:9003")
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
	v.SetDefault("ai.system_prompt", "你是旅行助手项目中的旅行智能体，只回答出行计划、旅游攻略、目的地、交通、酒店、天气和美食相关问题。")
	v.SetDefault("ai.max_prompt_chars", 2000)
	v.SetDefault("rag.enabled", true)
	v.SetDefault("rag.provider", "local")
	v.SetDefault("rag.address", "115.190.209.83:19530")
	v.SetDefault("rag.collection_name", "travel_knowledge")
	v.SetDefault("rag.embedding_dim", 768)
	v.SetDefault("rag.top_k", 3)
	v.SetDefault("rag.min_score", 0.15)
	v.SetDefault("travel_data.enabled", true)
	v.SetDefault("travel_data.amap_key", "")
	v.SetDefault("travel_data.qweather_key", "")
	v.SetDefault("travel_data.timeout", "5s")
}
