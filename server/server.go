package server

import (
	"github.com/gin-gonic/gin"

	"home/server/boiler"
)

func Run() {
	router := gin.Default()

	boilerRouter := router.Group("/boiler")
	{
		boilerRouter.GET("/state", boiler.GetState)
		boilerRouter.POST("/", boiler.GetState)
	}

	router.Run()
}
