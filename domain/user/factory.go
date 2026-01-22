package user

import (
	"errors"
)

type UserFactory struct{
	hasher PasswordHasher
	Validator EmailValidator
}

func NewUserFactory(hasher PasswordHasher, validatory EmailValidator) *UserFactory{
	return &UserFactory{
		hasher: hasher,
		Validator: validatory,
	}
}

func (f *UserFactory) CreateUser(id int64, firstName, lastName, email, password string, position UserType) (*User, error){
	return newUser(id, firstName, lastName, email, password, position, f.hasher, f.Validator)
}

func (f *UserFactory) CreateUnsafeUser(id int64, firstName, lastName, email, password string, position UserType) (*User, error){
	if email != "" && !f.Validator.ValidateFormat(email){
		return  nil, errors.New("Invalid Email format")
	}

	if firstName == ""{
		return nil, errors.New("user first name cannot be empty")
	}

	user := &User{
		ID: id,
		FirstName: firstName,
		LastName: lastName,
		Email: email,
		Position: position,
	}

	if password != ""{
		if err := user.SetPassword(password, f.hasher); err != nil{
			return nil, err
		}
	}

	return user, nil
}