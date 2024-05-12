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

	"github.com/atanda0x/aggregator/db/sqlc"
	"github.com/atanda0x/aggregator/handler"
	"github.com/atanda0x/aggregator/helper"
)

type apiConfig struct {
	DB *sqlc.Queries
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
		helper.ResWithError(c.Writer, http.StatusNotFound, fmt.Sprintf("Error parsing JSON: %v", err))
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
		helper.ResWithError(c.Writer, http.StatusNotFound, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	helper.ResWithJSON(c.Writer, http.StatusOK, user)
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
	router.POST("/user", apiCfg.CreateUserHandle)

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
