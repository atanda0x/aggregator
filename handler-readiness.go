package main

import (
	"net/http"

	"github.com/atanda0x/aggregator/helper"
	"github.com/gin-gonic/gin"
)

func HandlerReadiness(c *gin.Context) {
	helper.ResWithJON(c.Writer, http.StatusOK, struct{}{})
}
