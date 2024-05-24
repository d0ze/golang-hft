# High frequency trader

This repo contains a framework to trade algorithmically on a market exchange.

It is written in golang, with docker support to run it.

Its currently implemented to operate on kraken, but new brokers can be implemented to operate simultaneously on different markets
The implementation allows different strategies to be implemented and applied simultaneously


## Trend & Strategies

The application keeps tracks of the last 30 candles (interval can be configured). A __Trend__ structure is implemented to keep
track on the candles and give informations about the trend (twap, smap, etc). Each time a new candle ticks, its sent to a channel
which a goroutine consumes, checking if any new profitable operation can be performed. 

This goroutine calls a __IStrategy__ implementation, which exposes 2 methods

```
Open(entities.ITrend, candle entities.Candle, balance *entities.Balance) *entities.Order
Close(trend entities.ITrend, candle entities.Candle, positions []*entities.Position) *entities.Order
```

which will check if is profitable to open any new position, or close an open one. They must return an
order to be submitted, which the overlying goroutine will send to another channel, consumed from the order sender


## Available Indicators

- TWAP
- SMA
- BB
- SMI

## Available Markets

### Kraken

- XBTEUR
- XBTUSD 
- XBTUSDT
- ETHEUR 
- ETHUSD 
- LTCEUR 

## Configuration

The following environment variables can be set to configure the application

- LOG_LEVEL (trace, debug, info, warning, error, default=info)
- KRAKEN_API_KEY - kraken api key
- KRAKEN_SECRET - kraken api secret
- OHLC_INTERVALS - which timeframes (in minutes) to consider in the run (dash separated list, defined in minutes, default=1-60)
- OHLC_SIZE - how many candles to keep for every timeframe (default=60)
- STRATEGY - strategy to run
- STRATEGY_INTERVAL_CHECK - for which candles timeframe (in minutes) the strategy will check for open/close orders (default=1)"`
- MARKETS - which markets to consider (dash separated list, default=ETHEUR-XBTEUR)
