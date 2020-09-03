package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gtihub.com/hariolate/tonneau/service/models"
	"gtihub.com/hariolate/tonneau/shared"
	"net/http"
	"strconv"
)

func (s *Service) MakeIsAuthMiddleware() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		uidRaw := c.GetHeader("X-User-Id")
		uid64, err := strconv.ParseUint(uidRaw, 10, 32)
		uid := uint(uid64)
		if err != nil {
			c.Set("is_auth", false)
			return
		}

		tokenRaw := c.GetHeader("X-Auth-Token")
		token, err := uuid.Parse(tokenRaw)
		if err != nil {
			c.Set("is_auth", false)
			return
		}

		user, err := s.GetUserByToken(models.NewToken(uid, token))
		if err != nil {
			c.Set("is_auth", false)
			return
		}

		c.Set("is_auth", true)
		c.Set("user", user)
		c.Set("uid", uid)
		c.Set("token", token)

		c.Next()
	}
}

func (s *Service) MakeAuthRequiredMiddleware() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		s.MakeIsAuthMiddleware()(c)
		isAuth, _ := c.Get("is_auth")

		if !isAuth.(bool) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorMessageResponse("invalid login info or not login"))
			return
		}

		c.Next()
	}
}
