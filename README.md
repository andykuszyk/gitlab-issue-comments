# Gitlab Issue Comments
The aim of this project is to allow comments to be submitted from a form on a web page (such as a blog) and for these comments to be stored in Gitlab Issues as a back end.

I have seen similar implementations using GitHub issues, but my particular use case is for Gitlab.

## What's included
* A Go web service in `cmd/main.go` which acts as a proxy for the GitLab API.
* A sample UI in `www`.

## Running locally (TL;DR)
* Build the web service container: 'make build'
* Run the docker-compose: 'make up'
* Visit http://localhost:8080 and add comments.

## Running locally (the details)
* The `docker-compose.yml` assumes an environment variable is set in your environment called `GITLAB_TOKEN`, which should contain an `api` scoped personal access token.
* The UI is hard-coded to use a private project of mine, so you will need to substitute out the project id (`18159567`) to a value of your own in `www/index.html`.

## Suggested usage
* Run the container somewhere (built from `make build`) with the environment variable `GITLAB_TOKEN` set.
* Submit comments from a webpage, using similar logic to the example in `./www/index.html`
