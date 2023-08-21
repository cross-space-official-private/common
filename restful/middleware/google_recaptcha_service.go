package middleware

import (
	recaptcha "cloud.google.com/go/recaptchaenterprise/v2/apiv1"
	recaptchapb "cloud.google.com/go/recaptchaenterprise/v2/apiv1/recaptchaenterprisepb"
	"context"
	"fmt"
	"github.com/cross-space-official/common/configuration"
	"github.com/cross-space-official/common/logger"
)

type (
	GoogleRecaptchaConfig struct {
		Enabled        bool    `mapstructure:"enabled"`
		ProjectID      string  `mapstructure:"project-id"`
		SiteKey        string  `mapstructure:"site-key"`
		ScoreThreshold float32 `mapstructure:"score-threshold"`
	}

	RecaptchaService interface {
		AssessToken(c context.Context, token, expectedAction string) bool
		ShouldSkipByProfileID(c context.Context, profileID string) bool
	}

	RecaptchaProfileWhitelistRepository interface {
		AddWhitelist(c context.Context, profileID string) error
		ShouldSkipByProfileID(c context.Context, profileID string) bool
	}

	googleRecaptchaServiceImpl struct {
		config        GoogleRecaptchaConfig
		whitelistRepo RecaptchaProfileWhitelistRepository
	}
)

func (g *googleRecaptchaServiceImpl) ShouldSkipByProfileID(c context.Context, profileID string) bool {
	if !g.config.Enabled {
		return true
	}

	if len(profileID) == 0 || g.whitelistRepo == nil {
		return false
	}

	return g.whitelistRepo.ShouldSkipByProfileID(c, profileID)
}

func (g *googleRecaptchaServiceImpl) AssessToken(c context.Context, token, expectedAction string) bool {
	if !g.config.Enabled {
		return true
	}

	if len(token) == 0 {
		return false
	}

	client, err := recaptcha.NewClient(c)
	if err != nil {
		return true
	}
	defer client.Close()

	event := &recaptchapb.Event{
		Token:   token,
		SiteKey: g.config.SiteKey,
	}

	assessment := &recaptchapb.Assessment{
		Event: event,
	}

	request := &recaptchapb.CreateAssessmentRequest{
		Assessment: assessment,
		Parent:     fmt.Sprintf("projects/%s", g.config.ProjectID),
	}

	response, err := client.CreateAssessment(c, request)
	if err != nil {
		logger.GetLoggerEntry(c).Warn("The recaptcha assessment creating failed with %s", err.Error())
		return false
	}

	if !response.TokenProperties.Valid || response.TokenProperties.Action != expectedAction {
		logger.GetLoggerEntry(c).Warn("The recaptcha token was invalid or the action does not match")
		return false
	}

	logger.GetLoggerEntry(c).Info("The recaptcha token score: %f", response.RiskAnalysis.Score)
	if response.RiskAnalysis.Score < g.config.ScoreThreshold {
		logger.GetLoggerEntry(c).Warn("The recaptcha token was invalid due to low score: ", response.RiskAnalysis.Score)
		return false
	}

	// Handles score if necessary, refer to https://cloud.google.com/recaptcha-enterprise/docs/create-assessment
	return true
}

func NewGoogleRecaptchaService() RecaptchaService {
	config := &GoogleRecaptchaConfig{}
	configuration.BuildSkipErrors("data.google-recaptcha", config)

	return &googleRecaptchaServiceImpl{config: *config, whitelistRepo: nil}
}

func NewGoogleRecaptchaServiceWithWhitelist(whitelistRepo RecaptchaProfileWhitelistRepository) RecaptchaService {
	config := &GoogleRecaptchaConfig{}
	configuration.BuildSkipErrors("data.google-recaptcha", config)

	return &googleRecaptchaServiceImpl{config: *config, whitelistRepo: whitelistRepo}
}
