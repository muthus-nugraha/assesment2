package repository

import (
	"assignment2/app/models"
	"assignment2/app/resource"
	"assignment2/config"
	"fmt"
)

type OrderRepository interface {
	GetOrders() ([]models.Order, error, int64)
	NewOrder(Order *models.Order, updateData resource.Order) error
	GetOrderById(id uint, preload bool) (models.Order, error)
	DeleteOrder(id int) error
}

func NewOrderRepository() OrderRepository {
	return &dbConn{
		connection: config.Connect(),
	}
}

func (db *dbConn) NewOrder(Order *models.Order, createData resource.Order) error {
	if Order.ID != 0 {
		var existingItem []models.Item
		db.connection.Where("order_id = ?", Order.ID).Find(&existingItem)
		var existingItemCount int = len(existingItem)
		var newItemCount int = len(createData.Items)
		if newItemCount < existingItemCount {
			db.connection.Debug().Unscoped().Model(models.Item{}).Where("order_id = ?", Order.ID).
				Order("id asc").
				Limit(existingItemCount - newItemCount).
				Offset(newItemCount).
				Delete(&models.Item{})
		}
	}

	Order.CustomerName = createData.CustomerName
	err := db.connection.Save(Order).Error
	if err != nil {
		return err
	}

	for eachIndex, eachItem := range createData.Items {
		var item models.Item
		db.connection.Where("order_id = ?", Order.ID).
			Order("id asc").
			Limit(1).Offset(eachIndex).Find(&item)
		item.OrderID = Order.ID
		item.ItemCode = eachItem.ItemCode
		item.Description = eachItem.Description
		item.Quantity = eachItem.Quantity
		err := db.connection.Save(&item).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *dbConn) GetOrders() ([]models.Order, error, int64) {
	var Order []models.Order
	var count int64
	connection := db.connection.Model(&Order).Preload("Items").Find(&Order)
	err := connection.Error
	if err != nil {
		return Order, err, 0
	}
	db.connection.Model(Order).Count(&count)
	return Order, nil, count
}

func (db *dbConn) GetOrderById(id uint, preload bool) (models.Order, error) {
	var Order models.Order
	connection := db.connection
	fmt.Println("OrderId :", id)
	connection = connection.Where("id = ?", id)
	if preload {
		connection = connection.Preload("Items")
	}
	connection = connection.First(&Order)
	err := connection.Error
	if err != nil {
		return Order, err
	}
	return Order, err
}

func (db *dbConn) DeleteOrder(id int) error {
	var Order models.Order
	err := db.connection.Delete(&Order, id).Error
	if err != nil {
		return err
	}
	return nil
}
