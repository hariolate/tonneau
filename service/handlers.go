package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gtihub.com/hariolate/tonneau/service/models"
	"gtihub.com/hariolate/tonneau/shared"
	"net/http"
)

func (s *Service) LoginHandler(c *gin.Context) {
	var body models.Login

	_ = c.BindJSON(&body)
	err := body.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorMessageResponse(err.Error()))
		return
	}

	user, err := s.GetUserByEmail(body.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, shared.ErrorMessageResponse("invalid email or password"))
		return
	}

	pwd := []byte(body.Password)
	pwdMatched := shared.ComparePasswords(user.Password, pwd)

	if pwdMatched {
		tok := models.NewTokenFor(user)
		if err := s.InsertUserLoginToken(tok); err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, shared.SuccessDataResponse(gin.H{
			"token": tok.Token,
			"uid":   tok.UID,
		}))
		return
	}

	c.JSON(http.StatusUnauthorized, shared.ErrorMessageResponse("invalid email or password"))
}

func (s *Service) SignupHandler(c *gin.Context) {
	var body models.Signup

	_ = c.BindJSON(&body)
	err := body.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorMessageResponse(err.Error()))
		return
	}

	if _, err := s.GetUserByEmail(body.Email); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUser, err := s.CreateUser(body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, shared.ErrorMessageResponse(err.Error()))
				return
			}
			c.JSON(http.StatusOK, shared.SuccessDataResponse(gin.H{
				"uid":   newUser.ID,
				"email": newUser.Email,
			}))
			return
		}
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusConflict, shared.ErrorMessageResponse("user already exists."))
}

func (s *Service) LogoutHandler(c *gin.Context) {
	uid, _ := c.Get("uid")
	tok, _ := c.Get("token")

	if err := s.RemoveUserLoginToken(models.NewToken(uid.(uint), tok.(uuid.UUID))); err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorMessageResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessMessageResponse("logout"))
}

func (s *Service) UpdatePasswordHandler(c *gin.Context) {
	uid, _ := c.Get("uid")

	var body models.NewPassword
	_ = c.BindJSON(&body)
	if err := body.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorMessageResponse(err.Error()))
		return
	}

	if err := s.UpdateUserPassword(uid.(uint), body); err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorMessageResponse(err.Error()))
		return
	}
}

func (s *Service) UpdateEmailHandler(c *gin.Context) {
	uid, _ := c.Get("uid")

	var body models.NewEmail
	_ = c.BindJSON(&body)
	if err := body.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorMessageResponse(err.Error()))
		return
	}

	if err := s.UpdateUserEmail(uid.(uint), body); err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorMessageResponse(err.Error()))
		return
	}
}

func (s *Service) UpdateProfileHandler(c *gin.Context) {
	uid, _ := c.Get("uid")
	part := c.Param("part")

	var body models.Value
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorMessageResponse(err.Error()))
		return
	}

	switch part {
	case "alias":
		if err := s.UpdateUserProfileAlias(uid.(uint), body); err != nil {
			c.JSON(http.StatusBadRequest, shared.ErrorMessageResponse(err.Error()))
			return
		}
	case "picture":
		if err := s.UpdateUserProfilePicture(uid.(uint), body); err != nil {
			c.JSON(http.StatusBadRequest, shared.ErrorMessageResponse(err.Error()))
			return
		}
	default:
		c.JSON(http.StatusBadRequest, shared.ErrorMessageResponse("unsupported field"))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessMessageResponse("set"))
}

func (s *Service) GetProfileHandler(c *gin.Context) {
	uid, _ := c.Get("uid")

	profile, err := s.GetUserProfile(uid.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorMessageResponse("no profile"))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessDataResponse(
		gin.H{
			"profile": profile,
		},
	))
}

func (s *Service) GetProfilePartHandler(c *gin.Context) {
	uid, _ := c.Get("uid")
	part := c.Param("part")

	profile, err := s.GetUserProfile(uid.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorMessageResponse("no profile"))
		return
	}

	var value models.Value

	switch part {
	case "alias":
		value.Value = profile.Alias
	//case "matches":
	//	value.Value = profile.Matches
	case "picture":
		value.Value = profile.Picture
	case "trophies":
		value.Value = profile.Trophies
	default:
		c.JSON(http.StatusBadRequest, shared.ErrorMessageResponse("unsupported field"))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessDataResponse(gin.H{
		"value": value.Value,
	}))
}

func (s *Service) SetupUserHandlersFor(router gin.IRouter) {
	router.
		Use(s.MakeShowBodyMiddleWare()).
		POST("/user/login", s.LoginHandler). // /api/user/login
		POST("/user/signup", s.SignupHandler).
		POST("/user/logout", s.MakeAuthRequiredMiddleware(), s.LogoutHandler).
		GET("/user/logout", s.MakeAuthRequiredMiddleware(), s.LogoutHandler).
		GET("/user/profile", s.MakeAuthRequiredMiddleware(), func(context *gin.Context) {
			context.Redirect(http.StatusPermanentRedirect, "/user/profile/get")
		}).
		GET("/user/profile/get", s.MakeAuthRequiredMiddleware(), s.GetProfileHandler).
		GET("/user/profile/get/:part", s.MakeAuthRequiredMiddleware(), s.GetProfilePartHandler).
		POST("/user/profile/set/:part", s.MakeAuthRequiredMiddleware(), s.UpdateProfileHandler).
		POST("/user/update_password", s.MakeAuthRequiredMiddleware(), s.UpdatePasswordHandler).
		POST("/user/update_email", s.MakeAuthRequiredMiddleware(), s.UpdateEmailHandler)
}
