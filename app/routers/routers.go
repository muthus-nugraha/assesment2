package routers

import (
	_ "assignment2/docs"
	"log"

	"assignment2/app/handler"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init() {
	OrderHandler := handler.NewOrderHandler()
	r := gin.Default()

	r.POST("/Order", OrderHandler.NewOrder)
	r.PUT("/Order", OrderHandler.NewOrder)
	r.GET("/Order", OrderHandler.GetOrderList)
	r.GET("/Order/:order_id", OrderHandler.GetOrderDetail)
	r.DELETE("/Order/:order_id", OrderHandler.DeleteOrder)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	log.Println("=========== Server started ===========")
	r.Run(":1337")
}
