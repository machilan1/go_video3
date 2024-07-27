package store

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	db *sqlx.DB
}

type User struct {
	Id         int
	Email      string
	Name       string
	IsAdmin    bool `db:"is_admin"`
	LastOnline *time.Time
	Password   *string   `json:"-"`
	CanUpdate  bool      `db:"can_update"`
	CreatedAt  time.Time `db:"created_at"`
	IsActive   bool      `db:"is_active"`
}

func newUserStore(DB *sqlx.DB) *UserStore {

	s := UserStore{
		db: DB,
	}
	return &s

}

type CreateUserParam struct {
	Name     string
	Email    string
	Password string
}

type FindOneUserParam struct {
	UserId int
}

type FindOneWithEmailParam struct {
	Email string
}

func (s *UserStore) Create(p CreateUserParam) error {
	_, err := s.db.Exec("insert into app_user (name ,email ,password) values($1 , $2 ,$3)", p.Name, p.Email, p.Password)
	if err != nil {
		return err
	}
	return nil

}

func (s *UserStore) FindOne(p FindOneUserParam) (User, error) {

	var u User

	return u, nil
}

func (s *UserStore) FindOneWithEmail(e string) (User, error) {

	var u User

	err := s.db.Get(&u,
		`select 
			id, 
			name, 
			email, 
			password, 
			can_update, 
			created_at,
			is_admin,
			is_active
		from app_user
		where email = $1 and is_active =true
		limit 1
		`, e)
	if err != nil {
		return u, err
	}
	return u, err
}
