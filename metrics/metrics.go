package metrics

import (
	"github.com/cross-space-official-private/common/configuration"
	"github.com/newrelic/go-agent/v3/newrelic"
)

var app *newrelic.Application

type newrelicConfig struct {
	ServiceName string `mapstructure:"service-name"`
	LicenseKey  string `mapstructure:"license-key"`
	Enabled     bool   `mapstructure:"enabled"`
}

func GetApplication() *newrelic.Application {
	newrelicConfig := &newrelicConfig{}
	configuration.BuildSkipErrors("data.newrelic", newrelicConfig)

	app, _ = newrelic.NewApplication(
		newrelic.ConfigAppName(newrelicConfig.ServiceName),
		newrelic.ConfigLicense(newrelicConfig.LicenseKey),
		newrelic.ConfigEnabled(newrelicConfig.Enabled),
	)

	return app
}
