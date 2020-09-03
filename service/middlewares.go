package service

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gtihub.com/hariolate/tonneau/service/models"
	"gtihub.com/hariolate/tonneau/shared"
	"io"
	"log"
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

func (s *Service) MakeShowBodyMiddleWare() func(ctx *gin.Context) {
	return DebugLogger()
}

type MyReadCloser struct {
	rc io.ReadCloser
	w  io.Writer
}

func (rc *MyReadCloser) Read(p []byte) (n int, err error) {
	n, err = rc.rc.Read(p)
	log.Println("run here", n, err)
	if n > 0 {
		if n, err := rc.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return n, err
}

func (rc *MyReadCloser) Close() error {
	return rc.rc.Close()
}

func DebugLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.Request.RequestURI)
		if c.Request.Method == http.MethodPost {
			var buf bytes.Buffer
			newBody := &MyReadCloser{c.Request.Body, &buf}
			c.Request.Body = newBody
			c.Next()
			log.Println(buf.String())
		} else {
			c.Next()
		}
	}
}
