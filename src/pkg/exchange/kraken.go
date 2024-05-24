package exchange

import (
	"fmt"
	"time"

	"github.com/d0ze/golang-hft/src/internal"
	"github.com/d0ze/golang-hft/src/pkg/domain/entities"
	krakenapi "github.com/d0ze/kraken-go-api-client"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type KrakenPriceType string

type KrakenOrderSide string

type KrakenOrderStatus string

const (
	TYPE_MARKET            KrakenPriceType = "market"
	TYPE_LIMIT             KrakenPriceType = "limit"
	TYPE_TAKE_PROFIT       KrakenPriceType = "take-profit"
	TYPE_TAKE_PROFIT_LIMIT KrakenPriceType = "take-profit-limit"
	TYPE_STOP_LOSS         KrakenPriceType = "stop-loss"
	TYPE_STOP_LOSS_LIMIT   KrakenPriceType = "stop-loss-limit"
)

const (
	SIDE_BUY  KrakenOrderSide = "buy"
	SIDE_SELL KrakenOrderSide = "sell"
)

const (
	STATUS_PENDING   KrakenOrderStatus = "pending"
	STATUS_OPEN      KrakenOrderStatus = "open"
	STATUS_CLOSED    KrakenOrderStatus = "closed"
	STATUS_CANCELLED KrakenOrderStatus = "canceled"
	STATUS_EXPIRED   KrakenOrderStatus = "expired"
)

var leverages = map[internal.Market]decimal.Decimal{
	internal.ETHEUR:  decimal.NewFromInt(5),
	internal.XBTEUR:  decimal.NewFromInt(5),
	internal.XBTUSDT: decimal.NewFromInt(5),
	internal.XBTUSD:  decimal.NewFromInt(5),
	internal.ETHUSD:  decimal.NewFromInt(5),
	internal.LTCEUR:  decimal.NewFromInt(3),
}

type IKrakenCli interface {
	GetOHLC(pair internal.Market, interval int) ([]entities.Candle, error)
	PlaceOrder(order *entities.Order) error
	GetOrder(id string) (*entities.Order, error)
	GetBalance() (*entities.Balance, error)
	GetOpenPositions(market internal.Market) ([]*entities.Position, error)
	GetMarketsData(markets []internal.Market) (entities.IMarkets, error)
	GetLeverage(market internal.Market) decimal.Decimal
}

type krakenCli struct {
	cli *krakenapi.KrakenAPI
}

func InitClient() {
	cli := krakenapi.New(internal.Config.KrakenApiKey, internal.Config.KrakenSecret)
	KrakenCli = &krakenCli{cli: cli}
}

func (c *krakenCli) GetLeverage(market internal.Market) decimal.Decimal {
	return leverages[market]
}

// returns a list of candles for the given interval and pair
func (c *krakenCli) GetOHLC(pair internal.Market, interval int) ([]entities.Candle, error) {
	resp, err := c.cli.OHLCWithInterval(Pair(pair), fmt.Sprintf("%d", interval))
	if err != nil {
		return []entities.Candle{}, err
	}
	var res []entities.Candle
	for _, ohlc := range resp.OHLC {
		res = append(res, entities.Candle{
			Open:      decimal.NewFromFloat(ohlc.Open),
			High:      decimal.NewFromFloat(ohlc.High),
			Low:       decimal.NewFromFloat(ohlc.Low),
			Close:     decimal.NewFromFloat(ohlc.Close),
			Timestamp: ohlc.Time,
		})
	}
	return res, nil
}

// place order on kraken
// returns the remote id of the order
func (c *krakenCli) PlaceOrder(order *entities.Order) error {
	resp, err := c.cli.AddOrder(
		Pair(order.Market),
		string(Side(order.Side)),
		string(Type(order.PriceType)),
		order.InitialVolume.String(),
		map[string]string{
			"leverage": fmt.Sprintf("%d", order.Leverage),
			"price":    order.LimitPrice.String(),
		})
	logrus.Infof("response: %v", resp)
	if err != nil {
		return err
	}
	return nil
}

// gets an order from kraken given the
// client id
func (c *krakenCli) GetOrder(id string) (*entities.Order, error) {
	return &entities.Order{}, nil
}

func (c *krakenCli) GetMarketsData(markets []internal.Market) (entities.IMarkets, error) {
	resp, err := c.cli.AssetPairs()
	if err != nil {
		return entities.NewMarkets(), err
	}
	result := entities.NewMarkets()
	for _, market := range markets {
		var info krakenapi.AssetPairInfo
		switch market {
		case internal.XBTUSD:
			info = resp.XXBTZUSD
		// TODO asset pair info response doesnt carry
		//  xbtusdt, add it in the library
		case internal.XBTUSDT:
			info = resp.XXBTZUSD
		case internal.ETHEUR:
			info = resp.XETHZEUR
		case internal.ETHUSD:
			info = resp.XETHZUSD
		case internal.XBTEUR:
			info = resp.XXBTZEUR
		case internal.LTCEUR:
			info = resp.XLTCZEUR
		default:
			panic("unknown market")
		}
		pair := fmt.Sprintf("%s%s", info.Base, info.Quote)
		costMin, _ := decimal.NewFromString(info.CostMin)
		orderMin, _ := decimal.NewFromString(info.OrderMin)
		result.SetMetadata(
			IPair(pair),
			info.LotDecimals,
			ICurrency(info.Altname[0:3]),
			ICurrency(info.Altname[3:6]),
			costMin,
			orderMin,
		)

	}
	return result, nil
}

func (c *krakenCli) GetBalance() (*entities.Balance, error) {
	resp, err := c.cli.TradeBalance(map[string]string{})
	if err != nil {
		return &entities.Balance{}, err
	}
	return &entities.Balance{
		TradeBalance:  decimal.NewFromFloat(resp.TradeBalance),
		InitialMargin: decimal.NewFromFloat(resp.MarginOP),
		FreeMargin:    decimal.NewFromFloat(resp.FreeMargin),
		Equity:        decimal.NewFromFloat(resp.Equity),
		MarginLevel:   decimal.NewFromFloat(resp.MarginLevel),
	}, nil
}

func (c *krakenCli) GetOpenPositions(market internal.Market) ([]*entities.Position, error) {
	resp, err := c.cli.OpenPositions(map[string]string{
		"market":  Pair(market),
		"docalcs": "true",
	})
	if err != nil {
		return []*entities.Position{}, err
	}
	var res []*entities.Position
	for id, position := range *resp {
		res = append(res, &entities.Position{
			Id:        id,
			Size:      decimal.NewFromFloat(position.Volume),
			Side:      ISide(position.PositionType),
			Market:    IPair(position.Pair),
			Realized:  decimal.NewFromFloat(float64(position.Net)),
			Status:    internal.PositionStatus(position.Status),
			CreatedAt: time.UnixMilli(int64(position.TradeTime)),
			Cost:      decimal.NewFromFloat(position.Cost),
		})
	}
	return res, nil
}
func Pair(pair internal.Market) string {
	switch pair {
	case internal.XBTEUR:
		return krakenapi.XXBTZEUR
	case internal.XBTUSD:
		return krakenapi.XXBTZUSD
	case internal.XBTUSDT:
		return krakenapi.XBTUSDT
	case internal.ETHEUR:
		return krakenapi.XETHZEUR
	case internal.ETHUSD:
		return krakenapi.XETHZUSD
	case internal.LTCEUR:
		return krakenapi.XLTCZEUR
	default:
		panic("unknown market")
	}
}

func IPair(pair string) internal.Market {
	switch pair {
	case krakenapi.XXBTZEUR:
		return internal.XBTEUR
	case krakenapi.XXBTZUSD:
		return internal.XBTUSD
	case krakenapi.XBTUSDT:
		return internal.XBTUSDT
	case krakenapi.XETHZEUR:
		return internal.ETHEUR
	case krakenapi.XETHZUSD:
		return internal.ETHUSD
	case krakenapi.XLTCZEUR:
		return internal.LTCEUR
	default:
		panic("unknown market")
	}
}

func Currency(currency internal.Currency) string {
	switch currency {
	case internal.XBT:
		return "XBT"
	case internal.ETH:
		return "ETH"
	case internal.EUR:
		return "EUR"
	case internal.USD:
		return "USD"
	case internal.USDT:
		return "USDT"
	case internal.LTC:
		return "LTC"
	default:
		panic("unknown currency")
	}
}

func ICurrency(currency string) internal.Currency {
	switch currency {
	case "XBT":
		return internal.XBT
	case "ETH":
		return internal.ETH
	case "EUR":
		return internal.EUR
	case "USD":
		return internal.USD
	case "USDT":
		return internal.USDT
	case "LTC":
		return internal.LTC
	default:
		panic("unknown currency")
	}
}

func Side(side internal.OrderSide) KrakenOrderSide {
	switch side {
	case internal.BUY:
		return SIDE_BUY
	case internal.SELL:
		return SIDE_SELL
	default:
		panic("unknown side")
	}
}

func ISide(side string) internal.OrderSide {
	switch side {
	case string(SIDE_BUY):
		return internal.BUY
	case string(SIDE_SELL):
		return internal.SELL
	default:
		panic("unknown side")
	}
}

func Type(t internal.PriceType) KrakenPriceType {
	switch t {
	case internal.LIMIT:
		return TYPE_LIMIT
	case internal.MARKET:
		return TYPE_MARKET
	default:
		panic("unknown order type")
	}
}

var KrakenCli IKrakenCli
