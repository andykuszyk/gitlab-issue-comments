gofmt:
	gofmt -w ./

test: gofmt
	GITLAB_URL="http://localhost:8081" GITLAB_TOKEN="token" go test ./... -v -p 1

watch-test:
	find . | grep -v .git | grep -e 'go$$' | entr -c make test 

up:
	docker-compose up --force-recreate --build -d

watch-up:
	find . | grep -v .git | entr -c make up 

build: test
	docker build -t andykuszyk/gitlab-issue-comments .

publish:
	docker push andykuszyk/gitlab-issue-comments 
