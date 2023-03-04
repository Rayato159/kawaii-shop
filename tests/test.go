package kawaiitests

import (
	"encoding/json"

	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/pkg/databases"
	"github.com/jmoiron/sqlx"
)

type testConfig struct {
	cfg config.IConfig
}

type ITestConfig interface {
	GetDb() *sqlx.DB
	GetConfig() config.IConfig
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
func (cfg *testConfig) GetConfig() config.IConfig {
	return cfg.cfg
}
func (cfg *testConfig) SetJwtAccessExpires(t int) {
	cfg.cfg.Jwt().SetJwtAccessExpires(t)
}
func (cfg *testConfig) SetJwtRefreshExpires(t int) {
	cfg.cfg.Jwt().SetJwtRefreshExpires(t)
}

func ToJsonStringtify(obj any) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}
