package main

import (
	"github.com/andykuszyk/gitlab-issue-comments/internal/gic"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/topics/:topicName/comments", gic.GetComments)
	r.POST("/topics/:topicName/comments", gic.PostComments)
	r.Run()
}
