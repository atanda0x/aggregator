package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/atanda0x/aggregator/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found")
	}

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/healthz", handlers.HandlerReadiness)
	router.GET("/err", handlers.HandlerErr)

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
