package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/SaeedAlian/econest/api/utils"
)

type Config struct {
	Env                                   string
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
	MaxOrdersInPage                       int32
	MaxWalletTransactionsInPage           int32
	SMTPHost                              string
	SMTPPort                              string
	SMTPEmail                             string
	SMTPPassword                          string
	SMTPTestRecipientAddress              string
	WebsiteUrl                            string
	WebsiteName                           string
	ResetPasswordWebsitePageUrl           string
	EmailVerificationWebsitePageUrl       string
	UploadsRootDir                        string
	ShipmentPrice                         float64
	OrderFeeFactor                        float64
}

var Env = InitConfig()

func InitConfig() Config {
	env := os.Getenv("ENV")
	projectRoot := getProjectRoot()

	defaultFilename := ".env"
	defaultEnvPath := filepath.Join(projectRoot, defaultFilename)

	filename := defaultFilename
	envPath := defaultEnvPath

	if env != "" {
		customFilename := fmt.Sprintf(".env.%s", env)
		customEnvPath := filepath.Join(projectRoot, customFilename)

		customEnvExists, err := utils.PathExists(customEnvPath)
		if err != nil || !customEnvExists {
			log.Printf(
				"warning: couldn't load env file %s, switching to default...\n",
				customFilename,
			)

			env = "default"
		} else {
			filename = customFilename
			envPath = customEnvPath
		}
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("warning: couldn't load env file %s: %v\n", filename, err)
	}

	return Config{
		Env:                                   env,
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
		MaxWalletTransactionsInPage:           int32(20),
		MaxProductsInPage:                     int32(15),
		MaxProductTagsInPage:                  int32(20),
		MaxProductOffersInPage:                int32(15),
		MaxProductAttributesInPage:            int32(15),
		MaxProductCommentsInPage:              int32(5),
		MaxProductCategoriesInPage:            int32(15),
		MaxOrdersInPage:                       int32(10),
		SMTPHost:                              getEnv("SMTP_HOST", ""),
		SMTPPort:                              getEnv("SMTP_PORT", ""),
		SMTPEmail:                             getEnv("SMTP_MAIL", ""),
		SMTPPassword:                          getEnv("SMTP_PASS", ""),
		SMTPTestRecipientAddress:              getEnv("SMTP_TEST_RECIPIENT_ADDRESS", ""),
		WebsiteUrl:                            getEnv("WEBSITE_URL", "http://localhost:5173"),
		WebsiteName:                           getEnv("WEBSITE_NAME", "EcoNest"),
		ResetPasswordWebsitePageUrl: getEnv(
			"RESET_PASS_WEBSITE_PAGE_URL",
			"http://localhost:5173/reset-password",
		),
		EmailVerificationWebsitePageUrl: getEnv(
			"EMAIL_VERIFICATION_WEBSITE_PAGE_URL",
			"http://localhost:5173/email-verify",
		),
		UploadsRootDir: getEnv("UPLOADS_ROOT_DIR", "uploads"),
		ShipmentPrice:  getEnvAsFloat64("SHIPMENT_PRICE", 10.0),
		OrderFeeFactor: getEnvAsFloat64("ORDER_FEE_FACTOR", 0.05),
	}
}

func getEnv(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && len(val) > 0 {
		return val
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if val, ok := os.LookupEnv(key); ok && len(val) > 0 {
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fallback
		}

		return v
	}

	return fallback
}

func getEnvAsFloat64(key string, fallback float64) float64 {
	if val, ok := os.LookupEnv(key); ok && len(val) > 0 {
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fallback
		}

		return v
	}

	return fallback
}

func getProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	panic("could not find project root (go.mod)")
}
