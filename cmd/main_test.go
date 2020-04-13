package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/andykuszyk/gitlab-issue-comments/internal/gic"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

const (
	baseUrl = "http://0.0.0.0:8080"
)

var (
	gitlab *gitlabMock
)

func Test_PostComment(t *testing.T) {
	given, when, then := NewMainTest(t)

	given.a_valid_comment()

	when.i_post_the_comment()

	then.the_returned_status_code_should_be(http.StatusNoContent).and().
		the_comment_should_be_saved_in_gitlab()
}

func NewMainTest(t *testing.T) (*mainTest, *mainTest, *mainTest) {
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
}

func (m *mainTest) and() *mainTest {
	return m
}

func (m *mainTest) a_valid_comment() {
	m.comment = gic.Comment{
		Subject: "subject",
		Body:    "body",
	}
}

func (m *mainTest) i_post_the_comment() {
	b, err := json.Marshal(m.comment)
	require.NoError(m.t, err)
	response, err := http.Post(fmt.Sprintf("%s/topics/test-topic/comments/", baseUrl), "application/json", bytes.NewBuffer(b))
	require.NoError(m.t, err)
	m.response = response
}

func TestMain(m *testing.M) {
	go main()
	go runGitlabMock()
	os.Exit(m.Run())
}

type gitlabMock struct {
	Comments []gic.Comment
}

func runGitlabMock() {
	gitlab = &gitlabMock{
		Comments: []gic.Comment{},
	}
	gitlab.start()
}

func (g *gitlabMock) start() {

}
