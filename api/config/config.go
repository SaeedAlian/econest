package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host                                  string
	Port                                  string
	DBHost                                string
	DBUser                                string
	DBPassword                            string
	DBName                                string
	DBPort                                string
	KeyServerRedisAddr                    string
	AuthRedisAddr                         string
	RotateKeyDays                         int64
	AccessTokenExpirationInMin            float64
	RefreshTokenExpirationInMin           float64
	CSRFTokenExpirationInMin              float64
	ForgotPasswordTokenExpirationInMin    float64
	EmailVerificationTokenExpirationInMin float64
	MaxUsersInPage                        int32
	MaxProductsInPage                     int32
	MaxProductTagsInPage                  int32
	MaxProductOffersInPage                int32
	MaxProductAttributesInPage            int32
	MaxProductCommentsInPage              int32
	MaxProductCategoriesInPage            int32
	MaxStoresInPage                       int32
	SMTPHost                              string
	SMTPPort                              string
	SMTPEmail                             string
	SMTPPassword                          string
	WebsiteUrl                            string
	WebsiteName                           string
	ResetPasswordWebsitePageUrl           string
	EmailVerificationWebsitePageUrl       string
	UploadsRootDir                        string
}

var Env = InitConfig()

func InitConfig() Config {
	godotenv.Load()

	return Config{
		Host:                                  getEnv("HOST", "http://localhost"),
		Port:                                  getEnv("PORT", "8080"),
		DBUser:                                getEnv("DB_USER", "postgres"),
		DBHost:                                getEnv("DB_HOST", "localhost"),
		DBPassword:                            getEnv("DB_PASSWORD", "postgres"),
		DBName:                                getEnv("DB_NAME", "postgres"),
		DBPort:                                getEnv("DB_PORT", "5432"),
		KeyServerRedisAddr:                    getEnv("KEY_SERVER_REDIS_ADDRESS", "localhost:6379"),
		AuthRedisAddr:                         getEnv("AUTH_REDIS_ADDRESS", "localhost:6379"),
		RotateKeyDays:                         getEnvAsInt("ROTATE_KEY_DAYS", 2),
		AccessTokenExpirationInMin:            float64(15),
		RefreshTokenExpirationInMin:           float64(60 * 24 * 7),
		CSRFTokenExpirationInMin:              float64(30),
		ForgotPasswordTokenExpirationInMin:    float64(10),
		EmailVerificationTokenExpirationInMin: float64(10),
		MaxUsersInPage:                        int32(10),
		MaxStoresInPage:                       int32(5),
		MaxProductsInPage:                     int32(15),
		MaxProductTagsInPage:                  int32(20),
		MaxProductOffersInPage:                int32(15),
		MaxProductAttributesInPage:            int32(15),
		MaxProductCommentsInPage:              int32(5),
		MaxProductCategoriesInPage:            int32(15),
		SMTPHost:                              getEnv("SMTP_HOST", ""),
		SMTPPort:                              getEnv("SMTP_PORT", ""),
		SMTPEmail:                             getEnv("SMTP_MAIL", ""),
		SMTPPassword:                          getEnv("SMTP_PASS", ""),
		WebsiteUrl:                            getEnv("SMTP_PASS", "http://localhost:5173"),
		WebsiteName:                           getEnv("SMTP_PASS", "EcoNest"),
		ResetPasswordWebsitePageUrl: getEnv(
			"SMTP_PASS",
			"http://localhost:5173/reset-password",
		),
		EmailVerificationWebsitePageUrl: getEnv(
			"SMTP_PASS",
			"http://localhost:5173/email-verify",
		),
		UploadsRootDir: getEnv("UPLOADS_ROOT_DIR", "uploads"),
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
