package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
)

type MainConfig struct {
	Port    int    `toml:"port"`
	AppName string `toml:"appName"`
	Host    string `toml:"host"`
}

type EmailConfig struct {
	Authcode string `toml:"authcode"`
	Email    string `toml:"email" `
}

type RedisConfig struct {
	RedisPort     int    `toml:"port"`
	RedisDb       int    `toml:"db"`
	RedisHost     string `toml:"host"`
	RedisPassword string `toml:"password"`
}

type MysqlConfig struct {
	MysqlPort         int    `toml:"port"`
	MysqlHost         string `toml:"host"`
	MysqlUser         string `toml:"user"`
	MysqlPassword     string `toml:"password"`
	MysqlDatabaseName string `toml:"databaseName"`
	MysqlCharset      string `toml:"charset"`
}

type JwtConfig struct {
	ExpireDuration int    `toml:"expire_duration"`
	Issuer         string `toml:"issuer"`
	Subject        string `toml:"subject"`
	Key            string `toml:"key"`
}

type Rabbitmq struct {
	RabbitmqPort     int    `toml:"port"`
	RabbitmqHost     string `toml:"host"`
	RabbitmqUsername string `toml:"username"`
	RabbitmqPassword string `toml:"password"`
	RabbitmqVhost    string `toml:"vhost"`
}

type RagModelConfig struct {
	RagEmbeddingModel string `toml:"embeddingModel"`
	RagChatModelName  string `toml:"chatModelName"`
	RagDocDir         string `toml:"docDir"`
	RagBaseUrl        string `toml:"baseUrl"`
	RagDimension      int    `toml:"dimension"`
}

type VoiceServiceConfig struct {
	VoiceServiceApiKey    string `toml:"voiceServiceApiKey"`
	VoiceServiceSecretKey string `toml:"voiceServiceSecretKey"`
}

type Config struct {
	EmailConfig        `toml:"emailConfig"`
	RedisConfig        `toml:"redisConfig"`
	MysqlConfig        `toml:"mysqlConfig"`
	JwtConfig          `toml:"jwtConfig"`
	MainConfig         `toml:"mainConfig"`
	Rabbitmq           `toml:"rabbitmqConfig"`
	RagModelConfig     `toml:"ragModelConfig"`
	VoiceServiceConfig `toml:"voiceServiceConfig"`
}

type RedisKeyConfig struct {
	CaptchaPrefix   string
	IndexName       string
	IndexNamePrefix string
}

var DefaultRedisKeyConfig = RedisKeyConfig{
	CaptchaPrefix:   "captcha:%s",
	IndexName:       "rag_docs:%s:idx",
	IndexNamePrefix: "rag_docs:%s:",
}

var config *Config

// InitConfig 初始化项目配置
func InitConfig() error {
	config = new(Config)
	configPath := os.Getenv("GOPHERAI_CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.toml"
	}
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		return fmt.Errorf("load config %s: %w", configPath, err)
	}
	applyEnvironmentOverrides(config, os.Getenv)
	return nil
}

func GetConfig() *Config {
	if config == nil {
		if err := InitConfig(); err != nil {
			panic(err)
		}
	}
	return config
}

func applyEnvironmentOverrides(conf *Config, getenv func(string) string) {
	setString := func(key string, target *string) {
		if value := getenv(key); value != "" {
			*target = value
		}
	}
	setInt := func(key string, target *int) {
		if value := getenv(key); value != "" {
			if parsed, err := strconv.Atoi(value); err == nil {
				*target = parsed
			}
		}
	}

	setString("GOPHERAI_HTTP_HOST", &conf.Host)
	setInt("GOPHERAI_HTTP_PORT", &conf.Port)

	setString("GOPHERAI_MYSQL_HOST", &conf.MysqlHost)
	setInt("GOPHERAI_MYSQL_PORT", &conf.MysqlPort)
	setString("GOPHERAI_MYSQL_USER", &conf.MysqlUser)
	setString("GOPHERAI_MYSQL_PASSWORD", &conf.MysqlPassword)
	setString("GOPHERAI_MYSQL_DATABASE", &conf.MysqlDatabaseName)
	setString("GOPHERAI_MYSQL_CHARSET", &conf.MysqlCharset)

	setString("GOPHERAI_REDIS_HOST", &conf.RedisHost)
	setInt("GOPHERAI_REDIS_PORT", &conf.RedisPort)
	setString("GOPHERAI_REDIS_PASSWORD", &conf.RedisPassword)
	setInt("GOPHERAI_REDIS_DB", &conf.RedisDb)

	setString("GOPHERAI_RABBITMQ_HOST", &conf.RabbitmqHost)
	setInt("GOPHERAI_RABBITMQ_PORT", &conf.RabbitmqPort)
	setString("GOPHERAI_RABBITMQ_USERNAME", &conf.RabbitmqUsername)
	setString("GOPHERAI_RABBITMQ_PASSWORD", &conf.RabbitmqPassword)
	setString("GOPHERAI_RABBITMQ_VHOST", &conf.RabbitmqVhost)

	setInt("GOPHERAI_JWT_EXPIRE_DURATION", &conf.ExpireDuration)
	setString("GOPHERAI_JWT_ISSUER", &conf.Issuer)
	setString("GOPHERAI_JWT_SUBJECT", &conf.Subject)
	setString("GOPHERAI_JWT_KEY", &conf.Key)

	setString("GOPHERAI_RAG_EMBEDDING_MODEL", &conf.RagEmbeddingModel)
	setString("GOPHERAI_RAG_CHAT_MODEL", &conf.RagChatModelName)
	setString("GOPHERAI_RAG_DOC_DIR", &conf.RagDocDir)
	setString("GOPHERAI_RAG_BASE_URL", &conf.RagBaseUrl)
	setInt("GOPHERAI_RAG_DIMENSION", &conf.RagDimension)

	setString("GOPHERAI_EMAIL_ADDRESS", &conf.Email)
	setString("GOPHERAI_EMAIL_AUTHCODE", &conf.Authcode)
	setString("GOPHERAI_TTS_API_KEY", &conf.VoiceServiceApiKey)
	setString("GOPHERAI_TTS_SECRET_KEY", &conf.VoiceServiceSecretKey)
}
