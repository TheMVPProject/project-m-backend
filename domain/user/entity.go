package user

import (
	"errors"
)

type UserType int

const (
	UserManegers UserType = 1
	UserEmployee UserType = 2
)

type User struct{
	ID int64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	PasswordHash string `json:"-"`
	Position UserType `json:"position"`
}

func newUser(id int64, firstName, lastName, email, password string, position UserType, hasher PasswordHasher, validator EmailValidator) (*User, error){
	if !validator.ValidateFormat(email){
		return nil, errors.New("invalid email format")
	}

	user := &User{
		ID: id,
		FirstName: firstName,
		LastName: lastName,
		Email: email,
		Position: position,
	}

	if password != ""{
		if err := user.SetPassword(password, hasher); err != nil{
			return nil, err
		}
	}

	if err := user.Validate(validator); err != nil{
		return nil, err
	}

	return user, nil
}

func (u *User) SetPassword(plaintTextPassword string, hasher PasswordHasher) error{
	if len(plaintTextPassword) < 8{
		return  errors.New("password must be at least 8 carecters long")
	}

	hash, err := hasher.Hash(plaintTextPassword)
	if err != nil{
		return err
	}
	u.PasswordHash = hash
	return nil
}

func (u *User) CheckPassword(plaintTextPassword string, hasher PasswordHasher) error{
	return hasher.Check(u.PasswordHash, plaintTextPassword)
}


func (u *User) Validate(validator EmailValidator) error{
	if u.FirstName == ""{
		return  errors.New("user first name cannot be empty")
	}

	if u.Email == ""{
		return errors.New("user email cannot be empty")
	}

	if !validator.IsPersonalProvider(u.Email){
		return errors.New("email must be from a personal provider")
	}

	if u.PasswordHash == ""{
		return errors.New("password must not be empty")
	}

	return  nil
}

func (u *User) ValidateEmailFormat(validator EmailValidator) error{
	if !validator.ValidateFormat(u.Email){
		return errors.New("invalid email format")
	}

	return nil
}