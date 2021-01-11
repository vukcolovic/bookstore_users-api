package services

import (
	"UdemyApp/bookstore_users-api/domain/users"
	"UdemyApp/bookstore_users-api/utils/crypto_utils"
	"UdemyApp/bookstore_users-api/utils/date_utils"
	"UdemyApp/bookstore_utils-go/rest_errors"
)

var(
	UsersService usersServiceInterface = &userService{}
)

type userService struct {

}

type usersServiceInterface interface {
	GetUser(int64) (*users.User, rest_errors.RestErr)
	CreateUser(users.User) (*users.User, rest_errors.RestErr)
	UpdateUser(bool, users.User) (*users.User, rest_errors.RestErr)
	DeleteUser(int64) rest_errors.RestErr
	Search(string) (users.Users, rest_errors.RestErr)
	LoginUser(users.LoginRequest) (*users.User, rest_errors.RestErr)
}

func (s *userService) GetUser(userId int64) (*users.User, rest_errors.RestErr) {
	if userId <= 0 {
		return nil, rest_errors.NewBadRequestError("Invalid user id")
	}
	result := &users.User{Id: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *userService)  CreateUser(user users.User) (*users.User, rest_errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.Status = users.StatusActive
	user.DateCreated = date_utils.GetNowDbFormat()
	user.Password = crypto_utils.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userService)  UpdateUser(isPartial bool, user users.User) (*users.User, rest_errors.RestErr) {
	current := &users.User{Id: user.Id}
	err := current.Get()
	if err != nil {
		return nil, err
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}
		if user.LastName != "" {
			current.LastName = user.LastName
		}
		if user.Email != "" {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}


	if err := current.Update(); err != nil {
		return nil, err
	}

	return current, nil
}

func (s *userService)  DeleteUser(userId int64) rest_errors.RestErr {
	user := &users.User{Id: userId}
	return user.Delete()
}

func (s *userService)  Search(status  string) (users.Users, rest_errors.RestErr) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}

func (s *userService) LoginUser(request users.LoginRequest) (*users.User, rest_errors.RestErr) {
	dao := &users.User{
		Email: request.Email,
		Password: crypto_utils.GetMd5(request.Password),
	}
	if err := dao.FindByEmailAndPassword(); err != nil {
		return nil, err
	}

	return dao, nil
}