package auth

import (
	"github.com/machilan1/go_video/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userStore *store.UserStore
}

func NewAuthService(store *store.UserStore) *AuthService {
	s := AuthService{
		userStore: store,
	}
	return &s
}

type LoginParam struct {
	Email    string
	Password string
}

type SignUpParam struct {
	Email    string
	Password string
	Name     string
}

func (as AuthService) SignUp(p SignUpParam) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Password), 12)
	if err != nil {
		return err
	}

	err = as.userStore.Create(store.CreateUserParam{Email: p.Email, Name: p.Name, Password: string(hash)})
	if err != nil {
		return err
	}
	return nil

}

func (as AuthService) Login(p LoginParam) (store.User, error) {
	user, err := as.userStore.FindOneWithEmail(p.Email)
	if err != nil {
		return store.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(p.Password))
	if err != nil {
		return store.User{}, err
	}

	return user, nil
}
