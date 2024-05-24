package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env"
)

const (
	ServerAddress    = "localhost:8080"
	DBPath           = "postgresql://admin:12345@localhost:5432/loyalty?sslmode=disable"
	TokenExpire      = time.Hour * 24
	SecretKey        = "09d25e094faa6ca2556c818166b7a9563b93f7099f6f0f4caa6cf63b88e8d3e7"
	MigrateSourceURL = "file://internal/db/migrations"
)

type Config struct {
	ServerAddr       string `env:"RUN_ADDRESS"`
	DBPath           string `env:"DATABASE_URI"`
	AccrualSysAddr   string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	MigrateSourceURL string `env:"MIGRATE_SOURCE_URL"`

	TokenExpire time.Duration
	SecretKey   string
}

func Init() *Config {
	cfg := Config{
		TokenExpire: TokenExpire,
		SecretKey:   SecretKey,
	}

	regStringVar(&cfg.ServerAddr, "a", ServerAddress, "Server address")
	regStringVar(&cfg.DBPath, "d", DBPath, "Server db path")
	regStringVar(&cfg.AccrualSysAddr, "r", "", "Accrual system address")
	regStringVar(&cfg.MigrateSourceURL, "m", MigrateSourceURL, "Migrate source URL")

	flag.Parse()

	flagServerAddress := getStringFlag("a")
	flagDataBaseURI := getStringFlag("d")
	flagMigrateSourceURL := getStringFlag("m")

	err := env.Parse(&cfg)

	if err != nil {
		panic(err)
	}

	if flagServerAddress != ServerAddress {
		cfg.ServerAddr = flagServerAddress
	}
	if flagDataBaseURI != DBPath {
		cfg.DBPath = flagDataBaseURI
	}
	if flagMigrateSourceURL != MigrateSourceURL {
		cfg.DBPath = flagDataBaseURI
	}

	log.Printf("ServerAddr %s, DBPath %s, AccrualSysAddr %s", cfg.ServerAddr, cfg.DBPath, cfg.AccrualSysAddr)

	return &cfg
}

func regStringVar(p *string, name string, value string, usage string) {
	if flag.Lookup(name) == nil {
		flag.StringVar(p, name, value, usage)
	}
}

func getStringFlag(name string) string {
	return flag.Lookup(name).Value.(flag.Getter).Get().(string)
}
