package logger

import (
	appContext "app/context"
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	// ログレベル: panic,fatal,error,warn,info,debug,trace
	logLevel = "debug"

	// ログに出力するjsonのキー
	requestIdKey = "requestId"

	subIdKey = "subId"
)

func newLogger() *logrus.Logger {
	l := logrus.StandardLogger()
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	lv, _ := logrus.ParseLevel(logLevel)
	l.SetLevel(lv)
	return l
}

var logger = newLogger()

func Log(ctx context.Context) *logrus.Entry {
	requestId := appContext.GetRequestId(ctx)
	subId := appContext.GetSub(ctx)
	return logger.WithFields(map[string]interface{}{
		requestIdKey: requestId,
		subIdKey:     subId,
	})
}
