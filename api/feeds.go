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

func (apiCfg *apiConfig) handlerCreateFeed(c *gin.Context, user sqlc.User) {
	type param struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(c.Request.Body)

	params := param{}
	err := decoder.Decode(&params)
	if err != nil {
		helper.ResWithError(c.Writer, http.StatusForbidden, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(c.Request.Context(), sqlc.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		helper.ResWithError(c.Writer, http.StatusForbidden, fmt.Sprintf("Couldn't create feed: %v", err))
		return
	}

	helper.ResWithJSON(c.Writer, http.StatusCreated, feed)
}

func (apiCfg *apiConfig) handlerGetFeeds(c *gin.Context) {
	feeds, err := apiCfg.DB.GetFeeds(c.Request.Context())
	if err != nil {
		helper.ResWithError(c.Writer, http.StatusForbidden, fmt.Sprintf("Couldn't get feeds: %v", err))
		return
	}
	helper.ResWithJSON(c.Writer, http.StatusCreated, feeds)
}
