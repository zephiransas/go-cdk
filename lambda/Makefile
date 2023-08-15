#
# ローカル実行用 Makefile
#
# デプロイ時はCDKで以下モジュールを使用してビルドを行うため、このMakefileは参照されない
# https://docs.aws.amazon.com/cdk/api/v2/docs/aws-lambda-go-alpha-readme.html
#

.PHONY: build test clean fmt

.DEFAULT_GOAL := build

build:
	export GO111MODULE=on
	env GOOS=linux go build -o build/todos/add cmd/todos/add/main.go
	env GOOS=linux go build -o build/todos/delete cmd/todos/delete/main.go
	env GOOS=linux go build -o build/todos/done cmd/todos/done/main.go
	env GOOS=linux go build -o build/todos/list cmd/todos/list/main.go
	env GOOS=linux go build -o build/todos/show cmd/todos/show/main.go

	env GOOS=linux go build -o build/authorizer cmd/authorizer/main.go
	env GOOS=linux go build -o build/login cmd/login/main.go

test:
	go test -v -p 1 ./...

clean:
	rm -rf ./build/todos
	rm -rf ./build/authorizer
	rm -rf ./build/login
	go clean -testcache

fmt:
	go fmt ./...
