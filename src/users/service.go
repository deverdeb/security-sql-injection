package users

import (
    "fmt"
)

var Service = UsersService{
    dao: &Dao,
}

type UsersService struct {
    dao *UsersDao
}

func (service UsersService) Authentication(login, password string) (*User, error) {
    user, err := service.dao.FindByLoginAndPassword(login, password)
    if err != nil {
        return nil, fmt.Errorf("error during authentication - error: %w", err)
    }
    return user, nil
}

func (service UsersService) Create(user *User, password string) (*User, error) {
    user, err := service.dao.Create(user, password)
    if err != nil {
        return nil, fmt.Errorf("user creation failed - error: %w", err)
    }
    return user, nil
}

func (service UsersService) FindAll() ([]*User, error) {
    return service.dao.FindAll()
}

func (service UsersService) SeachByFirstnameOrLastname(search string) ([]*User, error) {
    return service.dao.SeachByFirstnameOrLastname(search)
}
