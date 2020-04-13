package gic

import (
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/swag"
	"github.com/xanzy/go-gitlab"
	"net/http"
	"os"
)

var (
	client *gitlab.Client
)

func init() {
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		panic("GITLAB_TOKEN is a required environment variable")
	}
	baseUrl := os.Getenv("GITLAB_URL")
	if baseUrl == "" {
		baseUrl = "https://api.gitlab.com"
	}
	gitlabClient, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseUrl))
	if err != nil {
		panic(err)
	}
	client = gitlabClient
}

func GetComments(c *gin.Context) {

}

func PostComments(c *gin.Context) {
	_, _, err := client.Issues.CreateIssue("", &gitlab.CreateIssueOptions{
		Title:       swag.String(""),
		Description: swag.String(""),
	})
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusNoContent)
}
