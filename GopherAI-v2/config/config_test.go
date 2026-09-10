package config

import "testing"

func TestApplyEnvironmentOverrides(t *testing.T) {
	conf := &Config{
		MainConfig:  MainConfig{Host: "127.0.0.1", Port: 9090},
		RedisConfig: RedisConfig{RedisHost: "127.0.0.1", RedisPort: 6379},
		MysqlConfig: MysqlConfig{MysqlHost: "127.0.0.1", MysqlPort: 3306},
		Rabbitmq:    Rabbitmq{RabbitmqHost: "127.0.0.1", RabbitmqPort: 5672},
	}

	applyEnvironmentOverrides(conf, func(key string) string {
		return map[string]string{
			"GOPHERAI_HTTP_PORT":         "19090",
			"GOPHERAI_MYSQL_HOST":        "mysql",
			"GOPHERAI_MYSQL_PASSWORD":    "mysql-password",
			"GOPHERAI_REDIS_HOST":        "redis",
			"GOPHERAI_RABBITMQ_HOST":     "rabbitmq",
			"GOPHERAI_JWT_KEY":           "jwt-secret",
			"GOPHERAI_RAG_EMBEDDING_MODEL": "embedding-model",
		}[key]
	})

	if conf.MainConfig.Port != 19090 || conf.MysqlHost != "mysql" || conf.RedisHost != "redis" || conf.RabbitmqHost != "rabbitmq" {
		t.Fatalf("expected service host and port overrides, got %+v", conf)
	}
	if conf.MysqlPassword != "mysql-password" || conf.JwtConfig.Key != "jwt-secret" || conf.RagEmbeddingModel != "embedding-model" {
		t.Fatalf("expected credential and model overrides, got %+v", conf)
	}
}
