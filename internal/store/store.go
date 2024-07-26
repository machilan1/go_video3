package store

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	DB           *sqlx.DB
	UserStore    *UserStore
	CourseStore  *CourseStore
	VideoStore   *VideoStore
	ChapterStore *ChapterStore
}

func NewStore(DBUrl string) (*Store, error) {

	db, err := sqlx.Connect("postgres", DBUrl)
	if err != nil {
		return nil, err
	}

	s := Store{
		DB:           db,
		UserStore:    newUserStore(db),
		CourseStore:  newCourseStore(db),
		VideoStore:   newVideoStore(db),
		ChapterStore: newChapterStore(db),
	}

	return &s, nil

}
