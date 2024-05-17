package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/atanda0x/aggregator/db/sqlc"
	"github.com/atanda0x/aggregator/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (apiCfG *apiConfig) handlerCreateFeedFollows(c *gin.Context, user sqlc.User) {
	type param struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(c.Request.Body)

	params := param{}
	err := decoder.Decode(&params)
	if err != nil {
		helper.ResWithError(c.Writer, http.StatusForbidden, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed_follows, err := apiCfG.DB.CreateFeedFollow(c.Request.Context(), sqlc.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		helper.ResWithError(c.Writer, http.StatusForbidden, fmt.Sprintf("Couggldn't create feed follow: %v", err))
		return
	}

	helper.ResWithJSON(c.Writer, http.StatusCreated, feed_follows)
}

func (apiCfg *apiConfig) handlerGetFeedFollows(c *gin.Context, user sqlc.User) {
	feed_follows, err := apiCfg.DB.GetFeedFollows(c.Request.Context(), user.ID)
	if err != nil {
		helper.ResWithError(c.Writer, http.StatusForbidden, fmt.Sprintf("Couldn't get feed follows: %v", err))
		return
	}

	helper.ResWithJSON(c.Writer, http.StatusCreated, feed_follows)
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(c *gin.Context, user sqlc.User) {
	feedFollowIDstr := c.Param("feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDstr)
	if err != nil {
		helper.ResWithError(c.Writer, http.StatusBadRequest, fmt.Sprintf("Couldn't parse feed follow id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(c.Request.Context(), sqlc.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})

	if err != nil {
		helper.ResWithError(c.Writer, http.StatusInternalServerError, fmt.Sprintf("Couldn't delete feed follow: %v", err))
		return
	}

	helper.ResWithJSON(c.Writer, http.StatusOK, struct{}{})
}
