package handler

import (
	"assignment2/app/helpers"
	"assignment2/app/models"
	"assignment2/app/repository"
	"assignment2/app/resource"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	repo repository.OrderRepository
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		repository.NewOrderRepository(),
	}
}

type OrderOut struct {
	ID           uint      `json:"order_id"`
	CustomerName string    `json:"customer_name"`
	OrderedAt    time.Time `json:"ordered_at"`
	Items        []ItemOut `gorm:"foreignKey:OrderID"`
}

type ItemOut struct {
	ItemID      uint   `json:"item_id"`
	ItemCode    string `json:"item_code"`
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
	OrderID     uint   `json:"order_id"`
}

func (h *OrderHandler) NewOrder(c *gin.Context) {
	repo := h.repo
	var req resource.Order
	err := c.ShouldBind(&req)
	if c.Request.Method == "PUT" && req.ID == 0 {
		response := helpers.APIResponse("bad request", http.StatusBadRequest, "error", "Please specify ID for update.")
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if err != nil {
		response := helpers.APIResponse("bad request", http.StatusBadRequest, "error", "Something went wrong IDK why.")
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	var Order models.Order
	if req.ID != 0 {
		Order, _ = repo.GetOrderById(req.ID, false)
		if Order.ID == 0 {
			response := helpers.APIResponse2("Order not found", http.StatusBadRequest, 0, 0, 0, "")
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		err := repo.NewOrder(&Order, req)
		if err != nil {
			response := helpers.APIResponse2("An error occured while trying to update order", 500, 0, 0, 0, "")
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
	} else {
		Order = models.Order{}
		err = repo.NewOrder(&Order, req)
		if err != nil {
			response := helpers.APIResponse2("An error occured while trying to update order", http.StatusBadRequest, 0, 0, 0, "")
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
	}

	orderResult, err := repo.GetOrderById(uint(Order.ID), true)
	filteredOrder := OrderOut{}
	filteredOrder.ID = orderResult.ID
	filteredOrder.CustomerName = orderResult.CustomerName
	filteredOrder.OrderedAt = orderResult.OrderedAt
	filteredOrder.Items = []ItemOut{}
	for _, itemValue := range orderResult.Items {
		eachItem := ItemOut{}
		eachItem.ItemID = *itemValue.ID
		eachItem.ItemCode = itemValue.ItemCode
		eachItem.Description = itemValue.Description
		eachItem.Quantity = itemValue.Quantity
		eachItem.OrderID = itemValue.OrderID

		filteredOrder.Items = append(filteredOrder.Items, eachItem)
	}

	response := helpers.APIResponse2("Success Create Order", http.StatusOK, 0, 0, 0, filteredOrder)
	c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) GetOrderList(c *gin.Context) {
	repo := h.repo
	result, err, count := repo.GetOrders()

	if err != nil {
		response := helpers.APIResponse("bad request", http.StatusBadRequest, "error", "Something went wrong IDK why.")
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	orderList := []OrderOut{}
	for _, value := range result {
		eachOrder := OrderOut{}
		eachOrder.ID = value.ID
		eachOrder.CustomerName = value.CustomerName
		eachOrder.OrderedAt = value.OrderedAt
		eachOrder.Items = []ItemOut{}
		for _, itemValue := range value.Items {
			eachItem := ItemOut{}
			eachItem.ItemID = *itemValue.ID
			eachItem.ItemCode = itemValue.ItemCode
			eachItem.Description = itemValue.Description
			eachItem.Quantity = itemValue.Quantity
			eachItem.OrderID = itemValue.OrderID

			eachOrder.Items = append(eachOrder.Items, eachItem)
		}
		orderList = append(orderList, eachOrder)
	}

	response := helpers.APIResponse2("Success Retreive Order", http.StatusOK, 0, 0, int(count), orderList)
	c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	repo := h.repo

	orderId, _ := strconv.Atoi(c.Param("order_id"))
	orderResult, err := repo.GetOrderById(uint(orderId), true)

	if err != nil {
		response := helpers.APIResponse2("Data not found", http.StatusBadRequest, 0, 0, 0, map[string]string{})
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	eachOrder := OrderOut{}
	eachOrder.ID = orderResult.ID
	eachOrder.CustomerName = orderResult.CustomerName
	eachOrder.OrderedAt = orderResult.OrderedAt
	eachOrder.Items = []ItemOut{}
	for _, itemValue := range orderResult.Items {
		eachItem := ItemOut{}
		eachItem.ItemID = *itemValue.ID
		eachItem.ItemCode = itemValue.ItemCode
		eachItem.Description = itemValue.Description
		eachItem.Quantity = itemValue.Quantity
		eachItem.OrderID = itemValue.OrderID

		eachOrder.Items = append(eachOrder.Items, eachItem)
	}

	response := helpers.APIResponse2("Data found", http.StatusOK, 0, 0, 0, eachOrder)
	c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("order_id"))
	if id == 0 {
		errorMessage := gin.H{"error_messages": "order_id can't be null"}
		response := helpers.APIResponse("bad request", http.StatusBadRequest, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	repo := h.repo
	delete := repo.DeleteOrder(id)
	if delete != nil {
		errorMessage := gin.H{"error_messages": "Something Error, please check the input"}
		response := helpers.APIResponse("bad request", http.StatusBadRequest, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	response := helpers.APIResponse2("Success Delete Order", http.StatusOK, 0, 0, 0, 0)
	c.JSON(http.StatusOK, response)
}
