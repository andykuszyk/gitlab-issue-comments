package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/andykuszyk/gitlab-issue-comments/internal/gic"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	baseUrl = "http://0.0.0.0:8080"
)

var (
	gitlab *gitlabMock
)

func init() {
	gitlab = &gitlabMock{
		Comments: []gic.Comment{},
	}
}

func Test_PostComment(t *testing.T) {
	given, when, then := NewMainTest(t)

	given.a_valid_comment()

	when.i_post_the_comment()

	then.the_returned_status_code_should_be(http.StatusNoContent).and().
		the_comment_should_be_saved_in_gitlab()
}

func Test_GetComments(t *testing.T) {
	given, when, then := NewMainTest(t)

	given.a_valid_comment().and().
		i_post_the_comment().and().
		i_post_the_comment()

	when.i_get_the_comments()

	then.the_returned_status_code_should_be(http.StatusOK).and().
		there_should_be_n_comments_returned(2)
}

func NewMainTest(t *testing.T) (*mainTest, *mainTest, *mainTest) {
	gitlab.Comments = []gic.Comment{}
	test := &mainTest{
		t: t,
	}
	return test, test, test
}

type mainTest struct {
	t        *testing.T
	comment  gic.Comment
	response *http.Response
}

func (m *mainTest) the_returned_status_code_should_be(code int) *mainTest {
	require.Equal(m.t, code, m.response.StatusCode)
	return m
}

func (m *mainTest) the_comment_should_be_saved_in_gitlab() {
	require.Len(m.t, gitlab.Comments, 1)
	require.NotNil(m.t, gitlab.Comments[0].CreatedAt)
	require.Equal(m.t, gitlab.Comments[0].Body, "body")
}

func (m *mainTest) and() *mainTest {
	return m
}

func (m *mainTest) a_valid_comment() *mainTest {
	currentTime := time.Now()
	m.comment = gic.Comment{
		CreatedAt: &currentTime,
		Body:      "body",
	}
	return m
}

func (m *mainTest) i_post_the_comment() *mainTest {
	b, err := json.Marshal(m.comment)
	require.NoError(m.t, err)
	response, err := http.Post(fmt.Sprintf("%s/topics/test-topic/comments/", baseUrl), "application/json", bytes.NewBuffer(b))
	require.NoError(m.t, err)
	m.response = response
	return m
}

func (m *mainTest) i_get_the_comments() *mainTest {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/topics/test-topic/comments", baseUrl), nil)
	require.NoError(m.t, err)
	request.Header.Add("content-type", "application/json")
	response, err := http.DefaultClient.Do(request)
	require.NoError(m.t, err)
	m.response = response
	return m
}

func (m *mainTest) there_should_be_n_comments_returned(n int) *mainTest {
	var comments []gic.Comment
	defer m.response.Body.Close()
	b, err := ioutil.ReadAll(m.response.Body)
	require.NoError(m.t, err)
	require.NotEmpty(m.t, string(b))
	err = json.Unmarshal(b, &comments)
	require.NoError(m.t, err)
	require.Len(m.t, comments, n)
	return m
}

func TestMain(m *testing.M) {
	go main()
	go gitlab.start()
	os.Exit(m.Run())
}

type gitlabMock struct {
	Comments []gic.Comment
}

func (g *gitlabMock) start() {
	r := gin.Default()
	r.POST("/api/v4/projects/:project/issues", g.handlePostIssues)
	r.GET("/api/v4/projects/:project/issues", g.handleGetIssues)
	r.Run(":8081")
}

var response = `
{
"project_id" : 4,
"id" : 84,
"created_at" : "2016-01-07T12:44:33.959Z",
"iid" : 14,
"title" : "Issues with auth",
"state" : "opened",
"assignees" : [],
"assignee" : null,
"labels" : [
"bug"
],
"upvotes": 4,
"downvotes": 0,
"merge_requests_count": 0,
"author" : {
"name" : "Alexandra Bashirian",
"avatar_url" : null,
"state" : "active",
"web_url" : "https://gitlab.example.com/eileen.lowe",
"id" : 18,
"username" : "eileen.lowe"
},
"description" : null,
"updated_at" : "2016-01-07T12:44:33.959Z",
"closed_at" : null,
"closed_by" : null,
"milestone" : null,
"subscribed" : true,
"user_notes_count": 0,
"due_date": null,
"web_url": "http://example.com/my-group/my-project/issues/14",
"references": {
"short": "#14",
"relative": "#14",
"full": "my-group/my-project#14"
},
"time_stats": {
"time_estimate": 0,
"total_time_spent": 0,
"human_time_estimate": null,
"human_total_time_spent": null
},
"confidential": false,
"discussion_locked": false,
"_links": {
"self": "http://example.com/api/v4/projects/1/issues/2",
"notes": "http://example.com/api/v4/projects/1/issues/2/notes",
"award_emoji": "http://example.com/api/v4/projects/1/issues/2/award_emoji",
"project": "http://example.com/api/v4/projects/1"
},
"task_completion_status":{
"count":0,
"completed_count":0
}
}`

type gitlabIssue struct {
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	Labels      *string    `json:"labels"`
}

func (g *gitlabMock) handlePostIssues(c *gin.Context) {
	log.Println("gitlab mock post received")
	issue := gitlabIssue{}
	err := c.BindJSON(&issue)
	if err != nil {
		bytes, e := ioutil.ReadAll(c.Request.Body)
		if e == nil {
			log.Println(string(bytes))
		}
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	c.Writer.Write([]byte(response))
	g.Comments = append(g.Comments, gic.Comment{
		Body:      issue.Description,
		CreatedAt: issue.CreatedAt,
	})
}

func (g *gitlabMock) handleGetIssues(c *gin.Context) {
	issues := "["
	for _, _ = range g.Comments {
		issues = fmt.Sprintf("%s%s,", issues, response)
	}
	issues = strings.TrimRight(issues, ",")
	issues = fmt.Sprintf("%s]", issues)
	c.Writer.Write([]byte(issues))
}
