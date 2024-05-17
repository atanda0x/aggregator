package api

import (
	"github.com/atanda0x/aggregator/db/sqlc"
	"github.com/atanda0x/aggregator/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type apiConfig struct {
	DB     *sqlc.Queries
	router *gin.Engine
}

func NewServer(q sqlc.Queries) *apiConfig {
	server := &apiConfig{DB: &q}

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/healthz", handler.HandlerReadiness)
	router.GET("/err", handler.HandlerErr)
	router.POST("/users", server.CreateUserHandle)
	router.GET("/users", server.middleWare(server.handlerGetUser))

	router.POST("/feeds", server.middleWare(server.handlerCreateFeed))
	router.GET("/feeds", server.handlerGetFeeds)

	router.POST("/feed_follows", server.middleWare(server.handlerCreateFeedFollows))
	router.GET("/feed_follows", server.middleWare(server.handlerGetFeedFollows))

	router.DELETE("/feed_follows/:feedFollowID", server.middleWare(server.handlerDeleteFeedFollow))

	router.GET("/posts", server.middleWare(server.handlerGetPostsForUser))

	server.router = router
	return server
}

func (s *apiConfig) Start(address string) error {
	return s.router.Run(address)
}
