package api

import (
	"fmt"
	"net/http"

	"github.com/atanda0x/aggregator/auth"
	"github.com/atanda0x/aggregator/db/sqlc"
	"github.com/atanda0x/aggregator/helper"
	"github.com/gin-gonic/gin"
)

type authHandler func(c *gin.Context, user sqlc.User)

func (apiConfig *apiConfig) middleWare(handler authHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey, err := auth.GetApiKey(c.Request.Header)
		if err != nil {
			helper.ResWithError(c.Writer, http.StatusForbidden, fmt.Sprintf("Auth error: %v", err))
			c.Abort()
			return
		}

		user, err := apiConfig.DB.GetUserByAPIKey(c.Request.Context(), apiKey)

		// user, err := apiCfg.DB.GetUserByAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			helper.ResWithError(c.Writer, http.StatusBadRequest, fmt.Sprintf("Couldn't get user: %v", err))
			c.Abort()
			return
		}

		handler(c, user)
	}
}
