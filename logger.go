package main

import "go.uber.org/zap"

var (
	logger *zap.SugaredLogger
)

func NewLogger() {
	plainlogger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	logger = plainlogger.Sugar()

}
