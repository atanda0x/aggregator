package handlers

import (
	"net/http"

	"github.com/atanda0x/aggregator/helper"
	"github.com/gin-gonic/gin"
)

func HandlerReadiness(c *gin.Context) {
	helper.ResWithJSON(c.Writer, http.StatusOK, struct{}{})
}

func HandlerErr(c *gin.Context) {
	helper.ResWithError(c.Writer, http.StatusOK, "Something went wrong!!!!")
}
