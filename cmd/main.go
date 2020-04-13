package main

import (
	"github.com/andykuszyk/gitlab-issue-comments/internal/gic"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/topics/:topicName/comments", gic.GetComments)
	r.POST("/topics/:topicName/comments", gic.PostComments)
	r.Run()
}
