package user


type PasswordHasher interface{
	Hash(password string) (string, error)
	Check(hashedPassword, password string) error
}

type EmailValidator interface{
	ValidateFormat(email string) bool
	IsPersonalProvider(domain string) bool
}