package boiler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"home/modules/boiler"
)

func GetState(c *gin.Context) {
	c.JSON(http.StatusOK, boiler.GetState())
}
