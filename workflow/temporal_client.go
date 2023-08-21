package workflow

import (
	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/client"
)

type TemporalConfig struct {
	url       string `mapstructure:"url"`
	namespace string `mapstructure:"namespace"`
}

type ClientFactory struct{}

var clientFactory ClientFactory

func Build(config TemporalConfig) client.Client {
	workflowClient, _ := client.Dial(client.Options{
		HostPort:  config.url,
		Namespace: config.namespace,
		Logger:    &temporalLogger{logger: logrus.New()},
	})

	return workflowClient
}

func GetTemporalClientFactory() *ClientFactory {
	return &clientFactory
}

type temporalLogger struct {
	logger *logrus.Logger
}

func (l *temporalLogger) Debug(msg string, keyvals ...interface{}) {
	l.logger.Debug(msg, keyvals)
}

func (l *temporalLogger) Info(msg string, keyvals ...interface{}) {
	l.logger.Info(msg, keyvals)
}

func (l *temporalLogger) Warn(msg string, keyvals ...interface{}) {
	l.logger.Warn(msg, keyvals)
}

func (l *temporalLogger) Error(msg string, keyvals ...interface{}) {
	l.logger.Error(msg, keyvals)
}
