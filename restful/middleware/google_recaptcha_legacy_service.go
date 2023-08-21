package middleware

import (
	"context"
	"github.com/cross-space-official-private/common/configuration"
	"github.com/cross-space-official-private/common/httpclient"
	"github.com/cross-space-official-private/common/logger"
	"time"
)

type (
	GoogleRecaptchaLegacyConfig struct {
		Enabled        bool    `mapstructure:"enabled"`
		SecretKey      string  `mapstructure:"secret-key"`
		ScoreThreshold float32 `mapstructure:"score-threshold"`
	}

	googleRecaptchaLegacyServiceImpl struct {
		config GoogleRecaptchaLegacyConfig
		client *httpclient.XSpaceHttpClient
	}

	reCHAPTCHAResponse struct {
		Success        bool      `json:"success"`
		ChallengeTS    time.Time `json:"challenge_ts"`
		Hostname       string    `json:"hostname,omitempty"`
		ApkPackageName string    `json:"apk_package_name,omitempty"`
		Action         string    `json:"action,omitempty"`
		Score          float32   `json:"score,omitempty"`
		ErrorCodes     []string  `json:"error-codes,omitempty"`
	}
)

func (g *googleRecaptchaLegacyServiceImpl) ShouldSkipByProfileID(c context.Context, profileID string) bool {
	return false
}

func (g *googleRecaptchaLegacyServiceImpl) AssessToken(c context.Context, token, expectedAction string) bool {
	if !g.config.Enabled {
		return true
	}

	res := reCHAPTCHAResponse{}

	response, err := g.client.BuildRequest(c).
		SetQueryParam("secret", g.config.SecretKey).
		SetQueryParam("response", token).
		SetResult(&res).
		Post("siteverify")

	logger.GetLoggerEntry(c).Info("The recaptcha challenge timestamp: ", res.ChallengeTS)
	if err != nil || response.IsError() {
		return false
	}

	if !res.Success {
		logger.GetLoggerEntry(c).Warn("The recaptcha token was invalid due to ", res.ErrorCodes)
		return false
	}

	logger.GetLoggerEntry(c).Info("The recaptcha token score: ", res.Score)
	if res.Score < g.config.ScoreThreshold {
		logger.GetLoggerEntry(c).Warn("The recaptcha token was invalid due to low score: ", res.Score)
		return false
	}

	return true
}

func NewGoogleRecaptchaLegacyService() RecaptchaService {
	config := &GoogleRecaptchaLegacyConfig{}
	configuration.BuildSkipErrors("data.google-recaptcha-legacy", config)

	client := httpclient.GetHttpClientFactory().Build("https://www.google.com/recaptcha/api/",
		httpclient.AuthPayload{}, nil)

	return &googleRecaptchaLegacyServiceImpl{config: *config, client: client}
}
