package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/atanda0x/aggregator/auth"
	"github.com/atanda0x/aggregator/db/sqlc"
	"github.com/atanda0x/aggregator/handler"
	"github.com/atanda0x/aggregator/helper"
)

type apiConfig struct {
	DB *sqlc.Queries
}

type authHandler func(c *gin.Context, user sqlc.User)

func (apiCfg *apiConfig) middleWare(handler authHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey, err := auth.GetApiKey(c.Request.Header)
		if err != nil {
			helper.ResWithError(c.Writer, http.StatusForbidden, fmt.Sprintf("Auth error: %v", err))
			c.Abort()
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			helper.ResWithError(c.Writer, http.StatusBadRequest, fmt.Sprintf("Couldn't get user: %v", err))
			c.Abort()
			return
		}

		handler(c, user)
	}
}

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

func (apiCfG *apiConfig) handlerCreateFeedFollow(c *gin.Context, user sqlc.User) {
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

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the evn")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to db:", err)
	}

	apiCfg := apiConfig{
		DB: sqlc.New(conn),
	}

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/healthz", handler.HandlerReadiness)
	router.GET("/err", handler.HandlerErr)
	router.POST("/users", apiCfg.CreateUserHandle)
	router.GET("/users", apiCfg.middleWare(apiCfg.handlerGetUser))

	router.POST("/feeds", apiCfg.middleWare(apiCfg.handlerCreateFeed))
	router.GET("/feeds", apiCfg.handlerGetFeeds)

	router.POST("/feed_follows", apiCfg.middleWare(apiCfg.handlerCreateFeedFollow))
	router.GET("/feed_follows", apiCfg.middleWare(apiCfg.handlerGetFeedFollows))

	router.DELETE("/feed_follows/:feedFollowID", apiCfg.middleWare(apiCfg.handlerDeleteFeedFollow))

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server starting on port %v", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
