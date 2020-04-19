package gic

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/swag"
	"github.com/xanzy/go-gitlab"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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
		baseUrl = "https://gitlab.com"
	}
	gitlabClient, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseUrl))
	if err != nil {
		panic(err)
	}
	client = gitlabClient
}

func GetComments(c *gin.Context) {
	topic := c.Param("topicName")
	issues, _, err := client.Issues.ListProjectIssues(topic, &gitlab.ListProjectIssuesOptions{
		Labels: gitlab.Labels{
			"gitlab-issue-comment",
		},
	})
	if err != nil {
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	comments := []Comment{}
	log.Printf("Found %d issues\n", len(issues))
	for _, issue := range issues {
		comments = append(comments, Comment{
			CreatedAt: issue.CreatedAt,
			Body:      issue.Description,
		})
	}
	bytes, err := json.Marshal(comments)
	log.Printf("Marshalled %d comments into the string: %s", len(comments), string(bytes))
	if err != nil {
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.Header().Add("content-type", "application/json")
	c.Writer.Write(bytes)
}

func PostComments(c *gin.Context) {
	comment := Comment{}
	err := c.BindJSON(&comment)
	if err != nil {
		log.Println(err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if comment.CreatedAt == nil {
		currentTime := time.Now()
		comment.CreatedAt = &currentTime
	}
	topic := c.Param("topicName")
	title := comment.Body
	if len(title) > 50 {
		title = title[:50]
	}
	issue, response, err := client.Issues.CreateIssue(topic, &gitlab.CreateIssueOptions{
		Title:       swag.String(title),
		CreatedAt:   comment.CreatedAt,
		Description: swag.String(comment.Body),
		Labels: &gitlab.Labels{
			"gitlab-issue-comment",
		},
	})
	if err != nil {
		log.Println(issue)
		bytes, e := ioutil.ReadAll(response.Body)
		if e == nil {
			log.Println(string(bytes))
		}
		log.Println(err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(fmt.Sprintf("%e", err)))
		return
	}
	c.Writer.WriteHeader(http.StatusNoContent)
}
