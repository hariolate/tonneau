package service

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gtihub.com/hariolate/tonneau/service/models"
)

func (s *Service) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	if err := s.db.Model(&user).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) InsertUserLoginToken(t *models.Token) error {
	return s.r.SAdd(s.c, t.RedisKey(), t.RedisToken()).Err()
}

func (s *Service) RemoveUserLoginToken(t *models.Token) error {
	return s.r.SRem(s.c, t.RedisKey(), t.RedisToken()).Err()
}

func (s *Service) CreateUser(u models.Signup) (*models.User, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	newUser := models.User{
		Email:    u.Email,
		Password: string(pass),
	}

	if err := s.db.Model(&models.User{}).Create(&newUser).Error; err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *Service) GetUserByToken(t *models.Token) (*models.User, error) {
	redisRes := s.r.SIsMember(s.c, t.RedisKey(), t.RedisToken())
	isTokenExist := redisRes.Err() == nil && redisRes.Val()
	if !isTokenExist {
		return nil, errors.New("login token not exists")
	}

	return s.GetUserByUID(t.UID)
}

func (s *Service) GetUserByUID(uid uint) (*models.User, error) {
	var user models.User
	if err := s.db.Model(&models.User{}).First(&user, uid).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) UpdateUserEmail(uid uint, newEmail models.NewEmail) error {
	user, err := s.GetUserByUID(uid)
	if err != nil {
		return err
	}

	user.Email = newEmail.Email
	if err := user.Validate(); err != nil {
		return err
	}

	return s.db.Model(&models.User{}).Save(user).Error
}

func (s *Service) UpdateUserPassword(uid uint, newPassword models.NewPassword) error {
	user, err := s.GetUserByUID(uid)
	if err != nil {
		return err
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(newPassword.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	user.Password = string(pass)
	if err := user.Validate(); err != nil {
		return err
	}

	if err := s.db.Model(&models.User{}).Save(user).Error; err != nil {
		return err
	}

	return s.InvalidateAllLoginTokens(uid)
}

func (s *Service) UpdateUserProfilePicture(uid uint, value models.Value) error {
	blob, ok := value.Value.(models.ProfilePicture)
	if !ok {
		return errors.New("invalid value")
	}
	user, err := s.GetUserByUID(uid)
	if err != nil {
		return err
	}

	user.Profile.Picture = blob

	return s.db.Model(&models.User{}).Save(user).Error
}

func (s *Service) UpdateUserProfileAlias(uid uint, value models.Value) error {
	alias, ok := value.Value.(models.ProfileAlias)
	if !ok {
		return errors.New("invalid value")
	}
	user, err := s.GetUserByUID(uid)
	if err != nil {
		return err
	}

	user.Profile.Alias = string(alias)

	return s.db.Model(&models.User{}).Save(user).Error
}

func (s *Service) RemoveUser(uid uint) error {
	user, err := s.GetUserByUID(uid)
	if err != nil {
		return err
	}

	if err := s.db.Model(&models.User{}).Delete(user).Error; err != nil {
		return err
	}

	redisTokens := fmt.Sprintf("user:%d:tokens", uid)
	return s.r.Del(s.c, redisTokens).Err()
}

func (s *Service) GetUserProfile(uid uint) (*models.Profile, error) {
	user, err := s.GetUserByUID(uid)
	if err != nil {
		return nil, err
	}

	return &user.Profile, nil
}

func (s *Service) AddMatchResultFor(uid uint, result models.MatchResult) error {
	user, err := s.GetUserByUID(uid)
	if err != nil {
		return err
	}

	isUserInMatch := false

	for _, player := range result.Players {
		if player.ID == user.ID {
			isUserInMatch = true
		}
	}

	if !isUserInMatch {
		return errors.New("user not in match")
	}

	if len(user.Profile.Matches) >= 20 {
		user.Profile.Matches = user.Profile.Matches[1:]
	}

	user.Profile.Matches = append(user.Profile.Matches, result)
	return nil
}

func (s *Service) InvalidateAllLoginTokens(uid uint) error {
	redisTokens := fmt.Sprintf("user:%d:tokens", uid)
	return s.r.Del(s.c, redisTokens).Err()
}
