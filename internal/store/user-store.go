package store

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	db *sqlx.DB
}

type User struct {
	Id          int `db:"id"`
	Email       string
	Name        string
	IsAdmin     bool       `db:"is_admin"`
	LastLoginAt *time.Time `db:"last_login_at"`
	Password    *string    `json:"-"`
	CanUpdate   bool       `db:"can_update"`
	CreatedAt   time.Time  `db:"created_at"`
	IsActive    bool       `db:"is_active"`
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

type FindUsersParams struct {
	Limit int
	Page  int
}

func (s *UserStore) Create(p CreateUserParam) error {
	_, err := s.db.Exec("insert into app_user (name ,email ,password) values($1 , $2 ,$3)", p.Name, p.Email, p.Password)
	if err != nil {
		return err
	}
	return nil

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

func (s *UserStore) FindMany(p FindUsersParams) ([]User, error) {

	users := []User{}

	index := 1

	clauses := []string{}
	values := []any{}

	selectClause := `
	select 
		id, 
		email, 
		name, 
		is_admin, 
		can_update ,
		last_login_at 
	from app_user`

	if p.Limit != 0 {
		clauses = append(clauses, " limit $$")
		values = append(values, p.Limit)
	}

	if offset := (p.Page - 1) * p.Limit; offset > 0 {
		clauses = append(clauses, " offset $$")
		values = append(values, offset)
	}

	for i, _ := range clauses {
		clauses[i] = strings.Replace(clauses[i], "$$", "$"+fmt.Sprintf("%d", index), 1)
		index++
	}

	limitClause := strings.Join(clauses, " ")

	final := selectClause + limitClause

	fmt.Println(final)

	err := s.db.Select(&users, final, values...)

	if err != nil {
		return nil, err
	}

	return users, nil
}
