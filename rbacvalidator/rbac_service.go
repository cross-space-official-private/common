package rbacvalidator

import (
	"context"
	"github.com/cross-space-official-private/common/httpclient"
	"github.com/cross-space-official-private/common/logger"
	"github.com/cross-space-official-private/common/restful"
	"github.com/gin-gonic/gin"
)

type (
	DomainObjectService interface {
		GetInstance(c *gin.Context) interface{}
	}
	ActionOwnerIDAccessor func(c *gin.Context) string
	ResourceIDAccessor    func(c *gin.Context) string
	ResourceIDsAccessor   func(c *gin.Context) []string

	CheckPermissionRequest struct {
		Resource  string                 `json:"resource"`
		Method    string                 `json:"method"`
		ActorID   string                 `json:"actor_id"`
		ObjectIDs map[string]interface{} `json:"object_ids"`
	}

	PermissionCheckBySpacesRequest struct {
		Requests []PermissionCheckSingleRequest `json:"requests"`
	}

	PermissionCheckSingleRequest struct {
		Method   string   `json:"method"`
		ActorID  string   `json:"actor_id"`
		SpaceIDs []string `json:"space_ids"`
	}

	RBACClientConfig interface {
		GetHost() string
		GetAPIKey() string
	}

	PermissionCheckBySpacesResponse struct {
		SpacePermissions []PermissionDTO `json:"space_permissions"`
	}

	PermissionDTO struct {
		Resource      string `json:"resource"`
		Method        string `json:"method"`
		ActorID       string `json:"actor_id"`
		SpaceID       string `json:"space_id"`
		HasPermission bool   `json:"has_permission"`
		Reason        string `json:"reason"`
	}
)

type RBACService interface {
	CheckPermission(c context.Context, req CheckPermissionRequest) bool
	// CheckPermissionForSpaces This method is used to check if the actor has permission to access the methods in multiple spaces
	CheckPermissionForSpaces(c context.Context, req PermissionCheckBySpacesRequest) map[string]map[string]map[string]bool
}

type RBACServiceImpl struct {
	client *httpclient.XSpaceHttpClient
	config RBACClientConfig
}

func NewRBACService(config RBACClientConfig) RBACService {
	return &RBACServiceImpl{
		client: httpclient.GetHttpClientFactory().Build(config.GetHost(), httpclient.AuthPayload{ApiKey: config.GetAPIKey()}, nil),
		config: config,
	}
}

func (s *RBACServiceImpl) CheckPermission(c context.Context, req CheckPermissionRequest) bool {
	result := &PermissionDTO{}
	errResult := restful.ErrorResponse{}
	response, err := s.client.BuildRequest(c).
		SetBody(req).
		SetResult(result).
		SetError(&errResult).
		Post("/internal/permission-check")
	if err != nil {
		logger.GetLoggerEntry(c).Errorf("Error from rbac-service %+v", err)
		return false
	} else if response.IsError() {
		logger.GetLoggerEntry(c).Errorf(errResult.Message)
		return false
	}

	logger.GetLoggerEntry(c).Infof("Validation result: %v Reason: %v", result.HasPermission, result.Reason)
	return result.HasPermission
}

func (s *RBACServiceImpl) CheckPermissionForSpaces(c context.Context, req PermissionCheckBySpacesRequest) map[string]map[string]map[string]bool {
	result := &PermissionCheckBySpacesResponse{}
	errResult := restful.ErrorResponse{}
	response, err := s.client.BuildRequest(c).
		SetBody(req).
		SetResult(result).
		SetError(&errResult).
		Post("/internal/permission-check/by-spaces")

	res := map[string]map[string]map[string]bool{}
	for _, r := range req.Requests {
		sub := map[string]map[string]bool{}
		for _, spaceID := range r.SpaceIDs {
			sub[spaceID][r.ActorID] = false
		}
		res[r.Method] = sub
	}

	if err != nil {
		logger.GetLoggerEntry(c).Errorf("Error from rbac-service %+v", err)
		return res
	} else if response.IsError() {
		logger.GetLoggerEntry(c).Errorf(errResult.Message)
		return res
	}

	for _, permission := range result.SpacePermissions {
		res[permission.Method][permission.SpaceID][permission.ActorID] = permission.HasPermission
	}

	return res
}
