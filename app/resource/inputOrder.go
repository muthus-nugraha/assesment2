package resource

import "time"

type Order struct {
	ID           uint             `json:"order_id" binding:"required"`
	CustomerName string           `json:"customer_name"`
	OrderAt      time.Time        `json:"order_at"`
	Items        []InputOrderItem `json:"items"`
}

type InputOrderItem struct {
	ItemID      uint   `json:"item_id"`
	ItemCode    string `json:"item_code"`
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
}
