package logger

import (
	"context"
	"github.com/sirupsen/logrus"
)

type loggerFactory struct {
	fields map[string]func(ctx context.Context) string
	logger *logrus.Logger
}

var factory = &loggerFactory{
	logger: buildLogger(),
	fields: make(map[string]func(ctx context.Context) string),
}

func WithField(key string, accessor func(ctx context.Context) string) {
	factory.fields[key] = accessor
}

func GetLoggerEntry(c context.Context) *logrus.Entry {
	return factory.build(c)
}

func buildLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger
}

func (lf *loggerFactory) build(c context.Context) *logrus.Entry {
	if c == nil {
		return logrus.NewEntry(lf.logger)
	}

	fields := make(logrus.Fields)
	for key, accessor := range lf.fields {
		fields[key] = accessor(c)
	}

	return lf.logger.WithFields(fields)
}
