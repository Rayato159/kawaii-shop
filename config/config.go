package config

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

type IConfig interface {
	App() IAppConfig
	Db() IDbConfig
	Jwt() IJwtConfig
}

type config struct {
	app *app
	db  *db
	jwt *jwt
}

type app struct {
	host         string
	port         uint
	name         string
	version      string
	readTimeout  time.Duration // Second
	writeTimeout time.Duration // Second
	bodyLimit    int           // Byte
	apiKey       string
	adminKey     string
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

type jwt struct {
	secretKey        string
	accessExpiresAt  int // Second
	refreshExpiresAt int // Second
}

type IAppConfig interface {
	Url() string
	Version() string
	Name() string
	ApiKey() string
	AdminKey() string
	BodyLimit() int
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
}

func (c *config) App() IAppConfig          { return c.app }
func (a *app) Url() string                 { return fmt.Sprintf("%s:%d", a.host, a.port) }
func (a *app) Version() string             { return a.version }
func (a *app) Name() string                { return a.name }
func (a *app) BodyLimit() int              { return a.bodyLimit }
func (a *app) ReadTimeout() time.Duration  { return a.readTimeout }
func (a *app) WriteTimeout() time.Duration { return a.writeTimeout }
func (a *app) ApiKey() string              { return a.apiKey }
func (a *app) AdminKey() string            { return a.adminKey }

type IDbConfig interface {
	Url() string
	MaxOpenConns() int
}

func (c *config) Db() IDbConfig { return c.db }
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
func (d *db) MaxOpenConns() int { return d.maxConnections }

type IJwtConfig interface {
	SecretKey() []byte
	AccessTokenExpires() int
	RefreshTokenExpires() int
	SetJwtAccessExpires(t int)
	SetJwtRefreshExpires(t int)
}

func (c *config) Jwt() IJwtConfig         { return c.jwt }
func (j *jwt) SecretKey() []byte          { return []byte(j.secretKey) }
func (j *jwt) AccessTokenExpires() int    { return j.accessExpiresAt }
func (j *jwt) RefreshTokenExpires() int   { return j.refreshExpiresAt }
func (j *jwt) SetJwtAccessExpires(t int)  { j.accessExpiresAt = t }
func (j *jwt) SetJwtRefreshExpires(t int) { j.refreshExpiresAt = t }

func LoadConfig(path string) IConfig {
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatalf("load dotenv failed: %v", err)
	}
	return &config{
		// App
		app: &app{
			host: envMap["APP_HOST"],
			port: func() uint {
				p, err := strconv.Atoi(envMap["APP_PORT"])
				if err != nil {
					log.Fatalf("app port is required")
				}
				return uint(p)
			}(),
			name:     envMap["APP_NAME"],
			version:  envMap["APP_VERSION"],
			apiKey:   envMap["APP_API_KEY"],
			adminKey: envMap["APP_ADMIN_KEY"],
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
		// Db
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
		// Jwt
		jwt: &jwt{
			secretKey: envMap["JWT_SECRET_KEY"],
			accessExpiresAt: func() int {
				exp, err := strconv.Atoi(envMap["JWT_ACCESS_EXPIRES"])
				if err != nil {
					log.Fatalf("access expires is required")
				}
				return exp
			}(),
			refreshExpiresAt: func() int {
				exp, err := strconv.Atoi(envMap["JWT_REFRESH_EXPIRES"])
				if err != nil {
					log.Fatalf("refresh expires is required")
				}
				return exp
			}(),
		},
	}
}
