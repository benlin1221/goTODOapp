package config

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memcache"
	"github.com/gofiber/storage/memory"
	"github.com/gofiber/storage/postgres"
	"github.com/gofiber/storage/redis"
	"github.com/gofiber/template/html"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Config struct {
	*viper.Viper
	errorHandler    fiber.ErrorHandler
	fiber           *fiber.Config
	database        *DatabaseConfig
	clientTLSConfig *tls.Config
}

var instantiated *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		instantiated = &Config{
			Viper: viper.New(),
		}

		// Set default configurations
		instantiated.setDefaults()

		// Select the .env file
		instantiated.SetConfigName(".env")
		instantiated.SetConfigType("dotenv")
		instantiated.AddConfigPath(".")

		// Automatically refresh environment variables
		instantiated.AutomaticEnv()

		// Read configuration
		if err := instantiated.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				fmt.Println("failed to read configuration:", err.Error())
				os.Exit(1)
			}
		}

		instantiated.setClientTLSConfig()

		instantiated.setFiberConfig()

		instantiated.setDatabaseConfig()
	})
	return instantiated
}

func (config *Config) setDefaults() {
	// Set default database configuration
	config.SetDefault("DB_DRIVER", "postgresql")
	config.SetDefault("DB_HOST", "localhost")
	config.SetDefault("DB_USERNAME", "")
	config.SetDefault("DB_PASSWORD", "")
	config.SetDefault("DB_PORT", 5432)
	config.SetDefault("DB_NAME", "maintainer")
}

func (config *Config) setClientTLSConfig() {
	clientCertFile := config.GetString("CLIENT_TLS_CERT")
	clientKeyFile := config.GetString("CLIENT_TLS_KEY")
	if clientCertFile != "" && clientKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		tlsCfg := tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12, // TLS versions below 1.2 are considered insecure - see https://www.rfc-editor.org/rfc/rfc7525.txt for details
		}
		if err != nil {
			log.Fatalf("Error creating x509 keypair from client cert file %s and client key file %s", clientCertFile, clientKeyFile)
		}
		config.clientTLSConfig = &tlsCfg
		return
	}
	config.clientTLSConfig = &tls.Config{
		MinVersion: tls.VersionTLS12, // TLS versions below 1.2 are considered insecure - see https://www.rfc-editor.org/rfc/rfc7525.txt for details
	}
}

func (config *Config) GetClientTLSConfig() *tls.Config {
	return config.clientTLSConfig
}

func (config *Config) getFiberViewsEngine() fiber.Views {
	var viewsEngine fiber.Views
	switch strings.ToLower(config.GetString("FIBER_VIEWS")) {
	// Use the official html/template package by default
	case "html":
		// Set file extension dynamically to FIBER_VIEWS
		if config.GetString("FIBER_VIEWS_EXTENSION") == "" {
			config.Set("FIBER_VIEWS_EXTENSION", ".html")
		}
		engine := html.New(config.GetString("FIBER_VIEWS_DIRECTORY"), config.GetString("FIBER_VIEWS_EXTENSION"))
		engine.AddFunc(
			"ToLower", func(s string) string {
				return strings.ToLower(s)
			},
		)
		engine.AddFunc(
			"ToTitle", func(s string) string {
				titleCaser := cases.Title(language.English)
				return titleCaser.String(s)
			},
		)
		engine.AddFunc(
			"MakeMap", func(values ...interface{}) (map[string]interface{}, error) {
				if len(values)%2 != 0 {
					return nil, errors.New("invalid dict call")
				}
				dict := make(map[string]interface{}, len(values)/2)
				for i := 0; i < len(values); i += 2 {
					key, ok := values[i].(string)
					if !ok {
						return nil, errors.New("dict keys must be strings")
					}
					dict[key] = values[i+1]
				}
				return dict, nil
			},
		)
		engine.AddFunc(
			"ToRFC3339", func(t time.Time) string {
				return t.Format(time.RFC3339)
			},
		)
		engine.Reload(config.GetBool("FIBER_VIEWS_RELOAD")).
			Debug(config.GetBool("FIBER_VIEWS_DEBUG")).
			Layout(config.GetString("FIBER_VIEWS_LAYOUT")).
			Delims(config.GetString("FIBER_VIEWS_DELIMS_L"), config.GetString("FIBER_VIEWS_DELIMS_R"))
		viewsEngine = engine
	default:
		panic("unsupported template engine")
	}
	return viewsEngine
}

func (config *Config) setDatabaseConfig() {
	config.database = &DatabaseConfig{
		Default: DatabaseDriver{
			Driver:   config.GetString("DB_DRIVER"),
			Host:     config.GetString("DB_HOST"),
			Username: config.GetString("DB_USERNAME"),
			Password: config.GetString("DB_PASSWORD"),
			DBName:   config.GetString("DB_NAME"),
			Port:     config.GetInt("DB_PORT"),
		},
	}
}

func (config *Config) getDatabaseConfig() *DatabaseConfig {
	return config.database
}

func (config *Config) setFiberConfig() {
	config.fiber = &fiber.Config{
		Prefork:                   config.GetBool("FIBER_PREFORK"),
		ServerHeader:              config.GetString("FIBER_SERVERHEADER"),
		StrictRouting:             config.GetBool("FIBER_STRICTROUTING"),
		CaseSensitive:             config.GetBool("FIBER_CASESENSITIVE"),
		Immutable:                 config.GetBool("FIBER_IMMUTABLE"),
		UnescapePath:              config.GetBool("FIBER_UNESCAPEPATH"),
		ETag:                      config.GetBool("FIBER_ETAG"),
		BodyLimit:                 config.GetInt("FIBER_BODYLIMIT"),
		Concurrency:               config.GetInt("FIBER_CONCURRENCY"),
		Views:                     config.getFiberViewsEngine(),
		ReadTimeout:               config.GetDuration("FIBER_READTIMEOUT"),
		WriteTimeout:              config.GetDuration("FIBER_WRITETIMEOUT"),
		IdleTimeout:               config.GetDuration("FIBER_IDLETIMEOUT"),
		ReadBufferSize:            config.GetInt("FIBER_READBUFFERSIZE"),
		WriteBufferSize:           config.GetInt("FIBER_WRITEBUFFERSIZE"),
		CompressedFileSuffix:      config.GetString("FIBER_COMPRESSEDFILESUFFIX"),
		ProxyHeader:               config.GetString("FIBER_PROXYHEADER"),
		GETOnly:                   config.GetBool("FIBER_GETONLY"),
		ErrorHandler:              config.errorHandler,
		DisableKeepalive:          config.GetBool("FIBER_DISABLEKEEPALIVE"),
		DisableDefaultDate:        config.GetBool("FIBER_DISABLEDEFAULTDATE"),
		DisableDefaultContentType: config.GetBool("FIBER_DISABLEDEFAULTCONTENTTYPE"),
		DisableHeaderNormalizing:  config.GetBool("FIBER_DISABLEHEADERNORMALIZING"),
		DisableStartupMessage:     config.GetBool("FIBER_DISABLESTARTUPMESSAGE"),
		ReduceMemoryUsage:         config.GetBool("FIBER_REDUCEMEMORYUSAGE"),
	}
}

func (config *Config) GetFiberConfig() *fiber.Config {
	return config.fiber
}

func (config *Config) GetSessionConfig() session.Config {
	var store fiber.Storage
	switch strings.ToLower(config.GetString("MW_FIBER_SESSION_STORAGE_PROVIDER")) {
	case "memory":
		sessionStore := memory.New(memory.Config{
			GCInterval: config.GetDuration("MW_FIBER_SESSION_STORAGE_GCINTERVAL"),
		})
		store = sessionStore
	case "memcache":
		sessionStore := memcache.New(memcache.Config{
			Servers: config.GetString("MW_FIBER_SESSION_STORAGE_HOST") + ":" + config.GetString("MW_FIBER_SESSION_STORAGE_PORT"),
			Reset:   config.GetBool("MW_FIBER_SESSION_STORAGE_RESET"),
		})
		store = sessionStore
	case "postgresql", "postgres":
		sessionStore := postgres.New(postgres.Config{
			Host:       config.GetString("MW_FIBER_SESSION_STORAGE_HOST"),
			Port:       config.GetInt("MW_FIBER_SESSION_STORAGE_PORT"),
			Username:   config.GetString("MW_FIBER_SESSION_STORAGE_USERNAME"),
			Password:   config.GetString("MW_FIBER_SESSION_STORAGE_PASSWORD"),
			Database:   config.GetString("MW_FIBER_SESSION_STORAGE_DATABASE"),
			Table:      config.GetString("MW_FIBER_SESSION_STORAGE_TABLE"),
			Reset:      config.GetBool("MW_FIBER_SESSION_STORAGE_RESET"),
			GCInterval: config.GetDuration("MW_FIBER_SESSION_STORAGE_GCINTERVAL"),
		})
		store = sessionStore
	case "redis":
		redisCfg := redis.Config{
			Host:     config.GetString("MW_FIBER_SESSION_STORAGE_HOST"),
			Port:     config.GetInt("MW_FIBER_SESSION_STORAGE_PORT"),
			Username: config.GetString("MW_FIBER_SESSION_STORAGE_USERNAME"),
			Password: config.GetString("MW_FIBER_SESSION_STORAGE_PASSWORD"),
			Database: config.GetInt("MW_FIBER_SESSION_STORAGE_DATABASE"),
			Reset:    config.GetBool("MW_FIBER_SESSION_STORAGE_RESET"),
		}
		if config.GetBool("MW_FIBER_SESSION_STORAGE_TLS_ENABLED") {
			redisCfg.TLSConfig = config.GetClientTLSConfig()
		}
		store = redis.New(redisCfg)
	}

	return session.Config{
		Expiration:     config.GetDuration("MW_FIBER_SESSION_EXPIRATION"),
		Storage:        store,
		KeyLookup:      fmt.Sprintf("cookie:%s", config.GetString("MW_FIBER_SESSION_COOKIENAME")),
		CookieDomain:   config.GetString("MW_FIBER_SESSION_COOKIEDOMAIN"),
		CookiePath:     config.GetString("MW_FIBER_SESSION_COOKIEPATH"),
		CookieSecure:   config.GetBool("MW_FIBER_SESSION_COOKIESECURE"),
		CookieHTTPOnly: config.GetBool("MW_FIBER_SESSION_COOKIEHTTPONLY"),
		CookieSameSite: config.GetString("MW_FIBER_SESSION_COOKIESAMESITE"),
	}
}
