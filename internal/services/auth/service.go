package auth

type Service interface {
	RegisterUser(name, email, password string) error
	Login(email, password string, appID uint) (string, error)
}
