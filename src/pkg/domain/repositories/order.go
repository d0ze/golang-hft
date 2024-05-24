package repositories

import (
	"context"

	"github.com/d0ze/golang-hft/src/internal"
	"github.com/d0ze/golang-hft/src/pkg/domain/entities"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IOrderRepository interface {
	Save(ctx context.Context, order *entities.Order)
	Find(ctx context.Context, id string) *entities.Order
	Update(ctx context.Context, order *entities.Order)
}

type orderRepository struct{}

func (r *orderRepository) Save(ctx context.Context, order *entities.Order) {
	logrus.Debugf("saving order %s", order.String())
	collection := internal.Database.Driver().Collection("orders")
	_, err := collection.InsertOne(ctx, &order)
	if err != nil {
		panic(err)
	}
}

func (r *orderRepository) Find(ctx context.Context, id string) *entities.Order {
	logrus.Debugf("retrieving order %s", id)
	collection := internal.Database.Driver().Collection("orders")
	filter := bson.D{{Key: "order_id", Value: id}}
	var result entities.Order
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		panic(err)
	}
	return &result
}

func (r *orderRepository) Update(ctx context.Context, order *entities.Order) {
	logrus.Debugf("updating order %s", order.Id)
	collection := internal.Database.Driver().Collection("orders")
	id, _ := primitive.ObjectIDFromHex(order.Id)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: order}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
}

func InitOrderRepository() {
	repo := &orderRepository{}
	OrderRepository = repo
}

var OrderRepository IOrderRepository
