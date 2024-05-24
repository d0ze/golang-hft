package goro

import (
	"github.com/d0ze/golang-hft/src/pkg/domain/entities"
	"github.com/d0ze/golang-hft/src/pkg/exchange"
	"github.com/sirupsen/logrus"
)

func HandleOrders(orders chan *entities.Order) {
	go func() {
		for order := range orders {
			logrus.Infof("handling order %v", order)
			err := exchange.KrakenCli.PlaceOrder(order)
			if err != nil {
				logrus.Warnf("error %v placing order", err)
			}
		}
	}()
}
