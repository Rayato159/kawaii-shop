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
