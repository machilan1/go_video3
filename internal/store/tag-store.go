package store

import "github.com/jmoiron/sqlx"

type TagStore struct {
	db    *sqlx.DB
	store *Store
}

type CourseTag struct {
	ID       int    `db:"id"`
	Label    string `db:"label"`
	Selected bool   `db:"-"`
}

func newTagStore(DB *sqlx.DB) *TagStore {
	s := TagStore{
		db: DB,
	}
	return &s
}

func (s *TagStore) FindTagsWithCourseID(ID int) ([]CourseTag, error) {
	var t []CourseTag
	err := s.db.Select(&t, `
		select 
			t.id,
			t.label
		from course_tag ct
		left join tag t on t.id = ct.tag_id
		where ct.course_id = $1
		`, ID)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TagStore) FindTags() ([]CourseTag, error) {
	var tags []CourseTag

	err := s.db.Select(&tags, "select label, id from tag")
	if err != nil {
		return nil, err
	}
	return tags, nil
}
