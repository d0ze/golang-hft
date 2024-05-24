package internal

import (
	"strings"

	"github.com/Netflix/go-env"
	log "github.com/sirupsen/logrus"
)

type config struct {
	LogLevel              string `env:"LOG_LEVEL,default=info"`
	DevelopmentMode       bool   `env:"DEVELOPMENT_MODE,default=false"`
	KrakenApiKey          string `env:"KRAKEN_API_KEY,default="`
	KrakenSecret          string `env:"KRAKEN_SECRET,default="`
	OHLCIntervals         string `env:"OHLC_INTERVALS,default=1-60"`
	OHLCSize              int    `env:"OHLC_SIZE,default=60"`
	Strategy              string `env:"STRATEGY,default=twap"`
	StrategyIntervalCheck int    `env:"STRATEGY_INTERVAL_CHECK,default=1"`
	Markets               string `env:"MARKETS,default=XBTEUR-ETHEUR"`
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

func (c *config) GetMarkets() []Market {
	markets := []Market{}
	configured := strings.Split(Config.Markets, "-")
	for _, market := range configured {
		markets = append(markets, IMarket(market))
	}
	return markets
}

func IMarket(market string) Market {
	switch market {
	case "XBTEUR":
		return XBTEUR
	case "XBTUSD":
		return XBTUSD
	case "XBTUSDT":
		return XBTUSDT
	case "ETHEUR":
		return ETHEUR
	case "ETHUSD":
		return ETHUSD
	case "LTCEUR":
		return LTCEUR
	default:
		panic("unknown market")
	}
}

var Config *config
