package gic

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/swag"
	"github.com/xanzy/go-gitlab"
	"log"
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
		baseUrl = "https://www.gitlab.com"
	}
	gitlabClient, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseUrl))
	if err != nil {
		panic(err)
	}
	client = gitlabClient
}

func GetComments(c *gin.Context) {
	topic := c.Param("topicName")
	issues, _, err := client.Issues.ListProjectIssues(topic, &gitlab.ListProjectIssuesOptions{})
	if err != nil {
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	comments := []Comment{}
	log.Printf("Found %d issues\n", len(issues))
	for _, issue := range issues {
		comments = append(comments, Comment{
			Subject: issue.Title,
			Body:    issue.Description,
		})
	}
	bytes, err := json.Marshal(comments)
	log.Printf("Marshalled %d comments into the string: %s", len(comments), string(bytes))
	if err != nil {
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.Write(bytes)
}

func PostComments(c *gin.Context) {
	topic := c.Param("topicName")
	_, _, err := client.Issues.CreateIssue(topic, &gitlab.CreateIssueOptions{
		Title:       swag.String(""),
		Description: swag.String(""),
	})
	if err != nil {
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(fmt.Sprintf("%e", err)))
		return
	}
	c.Writer.WriteHeader(http.StatusNoContent)
}
