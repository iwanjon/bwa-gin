package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUser) (User, error)
	LoginUser(input LoginInput) (User, error)
	CheckEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(id int, filelocation string) (User, error)
	GetUserById(id int) (User, error)
	GetAllUsers() ([]User, error)
	UodateUser(i FormUpdateInput) (User, error)
}

type service struct {
	repo Repository
}

func NewSevice(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUser) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"
	NewUser, err := s.repo.Save(user)
	if err != nil {
		return NewUser, err
	}
	return NewUser, nil
}

func (s *service) LoginUser(input LoginInput) (User, error) {

	email := input.Email
	password := input.Password
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("no user email found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *service) CheckEmailAvailable(input CheckEmailInput) (bool, error) {

	email := input.Email

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return false, err
	}
	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) SaveAvatar(id int, filelocation string) (User, error) {

	user, err := s.repo.FindById(id)
	if err != nil {
		return user, err
	}
	user.AvatarFileName = filelocation
	updateeduser, err := s.repo.Update(user)
	if err != nil {
		return updateeduser, err
	}
	return updateeduser, nil

}

func (s *service) GetUserById(id int) (User, error) {
	user, err := s.repo.FindById(id)
	if err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("no user found")
	}
	return user, nil

}

func (s *service) GetAllUsers() ([]User, error) {
	users, err := s.repo.FindAllUsers()
	if err != nil {
		return users, err
	}
	return users, nil
}

func (s *service) UodateUser(i FormUpdateInput) (User, error) {
	u, err := s.repo.FindById(i.ID)
	if err != nil {
		return u, err
	}
	u.Name = i.Name
	u.Occupation = i.Occupation
	u.Email = i.Email
	u, err = s.repo.Update(u)
	if err != nil {
		return u, err
	}
	return u, nil
}
