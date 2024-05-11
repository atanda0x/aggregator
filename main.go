package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/atanda0x/aggregator/helper"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func handlerReadiness(c *gin.Context) {
	helper.ResWithJSON(c.Writer, http.StatusOK, struct{}{})
}

func handlerErr(c *gin.Context) {
	helper.ResWithError(c.Writer, http.StatusOK, "Something went wrong!!!!")
}
func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found")
	}

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/healthz", handlerReadiness)
	router.GET("/err", handlerErr)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server starting on port %v", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
