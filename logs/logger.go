package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	App *zap.Logger
}

var Logger *logger

func init() {
	Logger = &logger{}
}

func (logger *logger) initializer(file string,level zapcore.Level,fileMaxSize int,maxFileCount int,maxBackupDays int,compress bool,serviceName string){
	logger.App = NewLogger(file, level, fileMaxSize, maxFileCount, maxBackupDays, compress, serviceName)
}
