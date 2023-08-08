#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

aws_local() {
  if hash awslocal 2>/dev/null; then
    awslocal "$@"
  else
    aws --endpoint "http://localhost:4566" "$@"
  fi
}

# dynamo db
aws_local dynamodb create-table \
  --cli-input-json "file://$SCRIPT_DIR/todos-go.json"