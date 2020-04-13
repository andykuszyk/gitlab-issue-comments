gofmt:
	gofmt -w ./

test: gofmt
	GITLAB_URL="http://localhost:8081" GITLAB_TOKEN="token" go test ./...

watch:
	find . | grep -v .git | grep -e 'go$$' | entr -c make test 
