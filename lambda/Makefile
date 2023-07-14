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
	env GOOS=linux go build -o build/list cmd/list/main.go

test:
	go test -v -p 1 ./...

clean:
	rm -rf ./build/list
	go clean -testcache

fmt:
	go fmt ./...
