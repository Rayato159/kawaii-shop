package config

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type IConfig interface {
	App() IAppConfig
	Db() IDbConfig
}

type config struct {
	app *app
	db  *db
}

type app struct {
	host         string
	port         uint
	version      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	bodyLimit    int
}

type db struct {
	host           string
	port           uint
	protocol       string
	username       string
	password       string
	database       string
	sslMode        string
	maxConnections int
}

type IAppConfig interface {
	Url() string
}

func (c *config) App() IAppConfig {
	return c.app
}
func (a *app) Url() string {
	return fmt.Sprintf("%s:%d", a.host, a.port)
}

type IDbConfig interface {
	Url() string
	MaxOpenConns() int
}

func (c *config) Db() IDbConfig {
	return c.db
}
func (d *db) Url() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.host,
		d.port,
		d.username,
		d.password,
		d.database,
		d.sslMode,
	)
}
func (d *db) MaxOpenConns() int {
	return d.maxConnections
}

func LoadConfig(path string) IConfig {
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatalf("load dotenv failed: %v", err)
	}
	return &config{
		app: &app{
			host: envMap["APP_HOST"],
			port: func() uint {
				p, err := strconv.Atoi(envMap["APP_PORT"])
				if err != nil {
					log.Fatalf("app port is required")
				}
				return uint(p)
			}(),
			version: envMap["APP_VERSION"],
			bodyLimit: func() int {
				s, err := strconv.Atoi(envMap["APP_BODY_LIMIT"])
				if err != nil {
					return 8 * 1024 * 1024 // 8 MiB
				}
				return s
			}(),
			readTimeout: func() time.Duration {
				t, err := strconv.Atoi(envMap["APP_READ_TIMEOUT"])
				if err != nil {
					return time.Second * 360
				}
				return time.Duration(int64(t) * int64(math.Pow10(9)))
			}(),
			writeTimeout: func() time.Duration {
				t, err := strconv.Atoi(envMap["APP_WRTIE_TIMEOUT"])
				if err != nil {
					return time.Second * 360
				}
				return time.Duration(int64(t) * int64(math.Pow10(9)))
			}(),
		},
		db: &db{
			host: envMap["DB_HOST"],
			port: func() uint {
				p, err := strconv.Atoi(envMap["DB_PORT"])
				if err != nil {
					log.Fatalf("db port is required")
				}
				return uint(p)
			}(),
			protocol: envMap["DB_PROTOCOL"],
			username: envMap["DB_USERNAME"],
			password: envMap["DB_PASSWORD"],
			database: envMap["DB_DATABASE"],
			sslMode:  envMap["DB_SSL_MODE"],
			maxConnections: func() int {
				m, err := strconv.Atoi(envMap["DB_MAX_CONNECTIONS"])
				if err != nil {
					log.Fatalf("db max connections is required")
				}
				return m
			}(),
		},
	}
}
