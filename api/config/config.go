package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host                        string
	Port                        string
	DBHost                      string
	DBUser                      string
	DBPassword                  string
	DBName                      string
	DBPort                      string
	KeyServerRedisAddr          string
	AuthRedisAddr               string
	RotateKeyDays               int64
	AccessTokenExpirationInMin  float64
	RefreshTokenExpirationInMin float64
	CSRFTokenExpirationInMin    float64
	MaxUsersInPage              int32
}

var Env = InitConfig()

func InitConfig() Config {
	godotenv.Load()

	return Config{
		Host:                        getEnv("HOST", "http://localhost"),
		Port:                        getEnv("PORT", "8080"),
		DBUser:                      getEnv("DB_USER", "postgres"),
		DBHost:                      getEnv("DB_HOST", "localhost"),
		DBPassword:                  getEnv("DB_PASSWORD", "postgres"),
		DBName:                      getEnv("DB_NAME", "postgres"),
		DBPort:                      getEnv("DB_PORT", "5432"),
		KeyServerRedisAddr:          getEnv("KEY_SERVER_REDIS_ADDRESS", "localhost:6379"),
		AuthRedisAddr:               getEnv("AUTH_REDIS_ADDRESS", "localhost:6379"),
		RotateKeyDays:               getEnvAsInt("ROTATE_KEY_DAYS", 2),
		AccessTokenExpirationInMin:  float64(15),
		RefreshTokenExpirationInMin: float64(60 * 24 * 7),
		CSRFTokenExpirationInMin:    float64(30),
		MaxUsersInPage:              int32(10),
	}
}

func getEnv(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fallback
		}

		return v
	}

	return fallback
}
