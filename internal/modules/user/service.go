package user

import (
	"errors"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with this email already exists")
)

type Service interface {
	CreateUser(req *CreateUserRequest) (*User, error)
	GetAllUsers() ([]User, error)
	GetUserByID(id uint) (*User, error)
	UpdateUser(id uint, req *UpdateUserRequest) (*User, error)
	DeleteUser(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateUser(req *CreateUserRequest) (*User, error) {
	existingUser, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	user := &User{
		Name:  req.Name,
		Email: req.Email,
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetAllUsers() ([]User, error) {
	return s.repo.FindAll()
}

func (s *service) GetUserByID(id uint) (*User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *service) UpdateUser(id uint, req *UpdateUserRequest) (*User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if req.Email != "" && req.Email != user.Email {
		existingUser, err := s.repo.FindByEmail(req.Email)
		if err != nil {
			return nil, err
		}
		if existingUser != nil {
			return nil, ErrUserAlreadyExists
		}
		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) DeleteUser(id uint) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return s.repo.Delete(id)
}
