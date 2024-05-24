package internal

import (
	"github.com/Netflix/go-env"
	log "github.com/sirupsen/logrus"
)

type config struct {
	LogLevel              string `env:"LOG_LEVEL,default=info"`
	DevelopmentMode       bool   `env:"DEVELOPMENT_MODE,default=false"`
	KrakenApiKey          string `env:"KRAKEN_API_KEY,default="`
	KrakenSecret          string `env:"KRAKEN_SECRET,default="`
	DbName                string `env:"DB_NAME,default=hft"`
	MigrationsPath        string `env:"MIGRATIONS_PATH,default=migrations/"`
	OHLCIntervals         string `env:"OHLC_INTERVALS,default=1#60"`
	OHLCSize              int    `env:"OHLC_SIZE,default=60"`
	Strategy              string `env:"STRATEGY,default=twap"`
	StrategyIntervalCheck int    `env:"STRATEGY_INTERVAL_CHECK,default=1"`
}

func (c *config) Parse() {
	_, err := env.UnmarshalFromEnviron(c)
	if err != nil {
		log.Fatalf("failed to parse common config: %v", err)
	}
}

func InitConfig() {
	conf := &config{}
	conf.Parse()
	Config = conf
}

var Config *config
