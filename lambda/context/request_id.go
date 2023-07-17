package context

import (
	"context"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

const requestIdKey = "REQUEST_ID"

func SetRequestId(ctx context.Context) context.Context {
	if lc, ok := lambdacontext.FromContext(ctx); ok {
		return context.WithValue(ctx, requestIdKey, lc.AwsRequestID)
	} else {
		return ctx // テスト時などローカル実行の場合は未設定
	}
}

func GetRequestId(ctx context.Context) string {
	v := ctx.Value(requestIdKey)
	id, ok := v.(string)
	if !ok {
		// ラムダ実行時に必須で設定される想定
		// ただし、ログ用途であることと、エラーチェックの煩雑さのトレードオフから、エラーにはしない
		return "none"
	}
	return id
}
