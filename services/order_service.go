package services

import (
	"seckill/datamodels"
	"seckill/repositories"
)

type IOrderService interface {
	GetOrderByID(int64) (*datamodels.Order, error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
	DeleteOrderByID(int64) bool
	InsertOrder(*datamodels.Order) (int64, error)
	UpdateOrder(*datamodels.Order) error
	InsertOrderByMessage(*datamodels.Message) (int64, error)
}

type OrderService struct {
	orderRepository repositories.IOrder
}

// 初始化函数
func NewOrderService(repository repositories.IOrder) IOrderService {
	return &OrderService{orderRepository : repository}
}

func (o *OrderService) GetOrderByID(orderID int64) (*datamodels.Order, error) {
	return o.orderRepository.SelectByKey(orderID)
}

func (o *OrderService) GetAllOrder() ([]*datamodels.Order, error) {
	return o.orderRepository.SelectAll()
}

func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	return o.orderRepository.SelectAllWithInfo()
}

func (o *OrderService) DeleteOrderByID(orderID int64) bool {
	return o.orderRepository.Delete(orderID)
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (int64, error) {
	return o.orderRepository.Insert(order)

}

func (o *OrderService) UpdateOrder(order *datamodels.Order) error {
	return o.orderRepository.Update(order)
}

func (o *OrderService) InsertOrderByMessage(message *datamodels.Message) (orderID int64, err error) {
	order := &datamodels.Order{
		UserId:      message.UserID,
		ProductId:   message.ProductID,
		OrderStatus: datamodels.OrderSuccess,
	}
	return o.InsertOrder(order)
}