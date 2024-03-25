package config

type ServerConfig struct {
	Port int `env:"SERVER_PORT"`
}

type RedisConfig struct {
	Password string `env:"REDIS_PASSWORD"`
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT"`
}

type CheckerConfig struct {
	N int `env:"N"`
	K int `env:"K"`
}
