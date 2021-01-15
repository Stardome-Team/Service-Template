package account

import (
	"github.com/Stardome-Team/Service-Template/pkg/logset"
)

type service struct {
	repo   Repository
	logger logset.Logger
}

// Service contains interfaces for authentication services
type Service interface {
	RegisterAccountWithProfile(request AuthenticationRequest) error
}

// NewService creates new authentication service
func NewService(r Repository, l logset.Logger) Service {
	return &service{repo: r, logger: l}
}

func (s *service) RegisterAccountWithProfile(request AuthenticationRequest) error {

	var act Account = Account{
		Username: request.UserName,
	}

	rowsAffected, err := s.repo.CreateAccount(act)

	if err != nil || rowsAffected == 0 {
		return nil
	}

	if err != nil {
		return nil
	}

	return nil
}
