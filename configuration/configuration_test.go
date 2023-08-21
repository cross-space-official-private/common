package configuration

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

type TestConfig struct {
	TestObj   TestObj `mapstructure:"test-obj"`
	TestValue int     `mapstructure:"test-value"`
}

type TestObj struct {
	FieldA string `mapstructure:"field-a"`
	FieldB string `mapstructure:"field-b"`
	FieldC string `mapstructure:"field-c"`
}

func TestConfigurationFactory(t *testing.T) {
	// Arrange
	Init("./")

	// Act
	config := &TestConfig{}
	subConfig := &TestObj{}
	Build("data", config)
	Build("data.test-obj", subConfig)

	// Assert
	assert.Equal(t, "1", config.TestObj.FieldA)
	assert.Equal(t, "2", config.TestObj.FieldB)
	assert.Equal(t, "3", config.TestObj.FieldC)
	assert.Equal(t, 8, config.TestValue)
	assert.Equal(t, config.TestObj, subConfig)
}
