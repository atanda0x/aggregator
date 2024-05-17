package api

import (
	"fmt"
	"net/http"

	"github.com/atanda0x/aggregator/db/sqlc"
	"github.com/atanda0x/aggregator/helper"
	"github.com/gin-gonic/gin"
)

func (apiCfg *apiConfig) handlerGetPostsForUser(c *gin.Context, user sqlc.User) {
	post, err := apiCfg.DB.GetPostsForUser(c.Request.Context(), sqlc.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		helper.ResWithError(c.Writer, http.StatusForbidden, fmt.Sprintf("Couldn't get poosts: %v", err))
		return
	}
	helper.ResWithJSON(c.Writer, http.StatusOK, post)
}
