package logger

import (
	appContext "app/context"
	"context"
	"github.com/sirupsen/logrus"
)

const (
	// ログレベル: panic,fatal,error,warn,info,debug,trace
	logLevel = "debug"

	// ログに出力するjsonのキー
	requestIdKey = "requestId"
)

func newLogger() *logrus.Logger {
	l := logrus.StandardLogger()
	l.SetFormatter(&logrus.JSONFormatter{})
	lv, _ := logrus.ParseLevel(logLevel)
	l.SetLevel(lv)
	return l
}

var logger = newLogger()

func Log(ctx context.Context) *logrus.Entry {
	requestId := appContext.GetRequestId(ctx)
	return logger.WithFields(map[string]interface{}{
		requestIdKey: requestId,
	})
}
