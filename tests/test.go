package kawaiitests

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/pkg/databases"
	"github.com/jmoiron/sqlx"
)

type testConfig struct {
	cfg config.IConfig
}

type ITestConfig interface {
	GetDb() *sqlx.DB
	GetJwtConfig() config.IJwtConfig
	SetJwtAccessExpires(t int)
	SetJwtRefreshExpires(t int)
}

func Setup() ITestConfig {
	cfg := config.LoadConfig("../../../.env.test")
	return &testConfig{
		cfg: cfg,
	}
}

func (cfg *testConfig) GetDb() *sqlx.DB {
	return databases.DbConnect(cfg.cfg.Db())
}
func (cfg *testConfig) GetJwtConfig() config.IJwtConfig {
	return cfg.cfg.Jwt()
}
func (cfg *testConfig) SetJwtAccessExpires(t int) {
	cfg.cfg.Jwt().SetJwtAccessExpires(t)
}
func (cfg *testConfig) SetJwtRefreshExpires(t int) {
	cfg.cfg.Jwt().SetJwtRefreshExpires(t)
}
