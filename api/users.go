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

func (apiCfg *apiConfig) CreateUserHandle(c *gin.Context) {
	type params struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(c.Request.Body)

	param := params{}
	err := decoder.Decode(&param)
	if err != nil {
		helper.ResWithError(c.Writer, http.StatusForbidden, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(c.Request.Context(), sqlc.CreateUserParams{
		ID:        uuid.New(),
		Name:      param.Name,
		Email:     param.Email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		helper.ResWithError(c.Writer, http.StatusBadRequest, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	helper.ResWithJSON(c.Writer, http.StatusOK, user)
}

func (apicfg *apiConfig) handlerGetUser(c *gin.Context, user sqlc.User) {
	helper.ResWithJSON(c.Writer, http.StatusOK, user)
}
