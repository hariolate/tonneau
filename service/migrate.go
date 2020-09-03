package service

import (
	"gtihub.com/hariolate/tonneau/service/models"
	"gtihub.com/hariolate/tonneau/shared"
)

func (s *Service) autoMigrate() {
	shared.NoError(s.db.AutoMigrate(&models.User{}, &models.Profile{}))
}
