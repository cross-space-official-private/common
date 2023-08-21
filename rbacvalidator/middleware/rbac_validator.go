package middleware

import (
	"fmt"
	"github.com/cross-space-official/common/consts"
	"github.com/cross-space-official/common/failure"
	"github.com/cross-space-official/common/rbacvalidator"
	"github.com/gin-gonic/gin"
)

const AuthKey = "has_authorized"

func SpaceRBACValidator(
	rbacService rbacvalidator.RBACService,
	method string,
	accessor rbacvalidator.ResourceIDAccessor,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		if val, exists := c.Get(AuthKey); exists && val == true {
			c.Next()
			return
		}

		userID := c.Request.Header.Get(consts.ProfileIDHeader)
		spaceID := accessor(c)

		request := rbacvalidator.CheckPermissionRequest{
			Resource: "space",
			Method:   method,
			ActorID:  userID,
			ObjectIDs: map[string]interface{}{
				"space_id": spaceID,
			},
		}

		res := rbacService.CheckPermission(c, request)
		if res {
			c.Next()
			return
		}

		_ = c.Error(failure.GeneratePlainFailure(failure.NoAuthorizationError, fmt.Sprintf("user has no authorization on method %s", method), ""))
		c.Abort()
	}
}
