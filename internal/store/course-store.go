package store

import (
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type CourseStore struct {
	db *sqlx.DB
}

type Course struct {
	ID          int         `db:"id"`
	Title       string      `db:"name"`
	Description string      `db:"description"`
	Instructor  string      `db:"instructor_name"`
	Views       int         `db:"click_count"`
	Tags        []CourseTag `db:"-"`
	Chapters    []Chapter   `db:"-"`
	CreatedBy   int         `db:"created_by"`
	UpdatedAt   time.Time   `db:"updated_at"`
}

type FindCoursesParams struct {
	Page   int
	Limit  int
	UserID int
}

type CreateCourseBody struct {
	Title       string
	Instructor  string
	Description string
	CreatedBy   int
	Tags        []int
}

type UpdateCourseBody struct {
	Title       string
	Instructor  string
	Tags        []int
	Description string
}

func newCourseStore(DB *sqlx.DB) *CourseStore {
	s := CourseStore{
		db: DB,
	}
	return &s
}

func (s *CourseStore) CreateCourse(b CreateCourseBody) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	var lastInsertId int

	err = tx.QueryRow("insert into course (name,instructor_name,description, created_by) values($1,$2,$3,$4) returning id", b.Title, b.Instructor, b.Description, b.CreatedBy).Scan(&lastInsertId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// insert tags

	var valuesPairs []string

	for _, v := range b.Tags {
		valuesPairs = append(valuesPairs, `(`+strconv.Itoa(v)+`,`+strconv.Itoa(lastInsertId)+`)`)
	}

	valuesString := strings.Join(valuesPairs, ",")

	if valuesString != "" {
		_, err = tx.Exec("insert into course_tag (tag_id,course_id) values " + valuesString)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (s *CourseStore) UpdateCourse(courseID int, b UpdateCourseBody) error {

	tr, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tr.Exec("delete from course_tag where course_id = $1", courseID)
	if err != nil {
		return err
	}

	var valuesPairs []string

	for _, v := range b.Tags {
		valuesPairs = append(valuesPairs, `(`+strconv.Itoa(v)+`,`+strconv.Itoa(courseID)+`)`)

	}

	valuesString := strings.Join(valuesPairs, ",")

	if valuesString != "" {
		_, err = tr.Exec("insert into course_tag (tag_id,course_id) values " + valuesString)
		if err != nil {
			return err
		}
	}
	//
	interpolationIndex := 1
	fieldListTemp := []string{}
	fieldList := []string{}
	valueList := []any{}

	if b.Description != "" {
		fieldListTemp = append(fieldListTemp, "description=")
		valueList = append(valueList, b.Description)
	}

	if b.Title != "" {
		fieldListTemp = append(fieldListTemp, "name=")
		valueList = append(valueList, b.Title)
	}
	if b.Instructor != "" {
		fieldListTemp = append(fieldListTemp, "instructor_name=")
		valueList = append(valueList, b.Instructor)
	}

	for _, v := range fieldListTemp {
		index := strconv.Itoa(interpolationIndex)
		fieldList = append(fieldList, v+`$`+index)
		interpolationIndex++
	}

	fieldsString := strings.Join(fieldList, ",")
	valueList = append(valueList, courseID)
	index := strconv.Itoa(interpolationIndex)

	_, err = tr.Exec("update course set "+fieldsString+` where id = $`+index, valueList...)

	if err != nil {
		return err
	}
	err = tr.Commit()
	if err != nil {
		return err
	}

	return err
}

// TODO : Figure out a way to refactor this
func (s *CourseStore) FindOne(ID int) (Course, error) {

	// find related chapters
	chapters := []Chapter{}
	course := Course{}

	err := s.db.Select(&chapters, `
	select 
		c.id, 
		c.title, 
		c.description, 
		c.chap_num,
		c.created_at,
		c.course_id,
		v.id as "video_id",
		v.file_name,
		'/statics/videos/'||v.file_name as "link"
	from chapter c
	join video v on v.id = c.video_id 
	where c.course_id = $1
	order by c.chap_num asc
		`, ID)

	if err != nil {
		return Course{}, err
	}

	err = s.db.Get(&course, `
	select 
		c.id,
		c.name,
		c.description, 
		c.instructor_name, 
		c.click_count ,
		c.created_by
		from course c 
	where c.id = $1
	`, ID)

	if err != nil {
		return Course{}, err
	}

	course.Chapters = chapters

	var t []CourseTag
	err = s.db.Select(&t, `
		select 
			t.id,
			t.label
		from course_tag ct
		left join tag t on t.id = ct.tag_id
		where ct.course_id = $1
		`, ID)
	course.Tags = t
	if err != nil {
		return Course{}, err
	}

	return course, nil

}

func (s *CourseStore) FindMany(p FindCoursesParams) ([]Course, error) {
	c := []Course{}
	courses := &c
	var err error

	if p.UserID != 0 {
		err = s.db.Select(courses, `
		select  
			c.id,
			c.name,
			c.description, 
			c.instructor_name, 
			c.click_count,
			c.created_by,
			c.updated_at
		from course c
		where c.created_by =$1
		order by c.updated_at desc
		offset $2 limit $3;
		`, p.UserID, (p.Page-1)*p.Limit, p.Limit)
		if err != nil {
			return nil, err
		}
	} else {
		err = s.db.Select(courses, `
	select  
		c.id,
		c.name,
		c.description, 
		c.instructor_name, 
		c.click_count,
		c.created_by,
		c.updated_at
	from course c
	order by c.updated_at desc
	offset $1 limit $2;
	`, (p.Page-1)*p.Limit, p.Limit)
		if err != nil {
			return nil, err
		}
	}

	for i, v := range *courses {
		var t []CourseTag
		err := s.db.Select(&t, `
		select 
			t.id,
			t.label
		from course_tag ct
		left join tag t on t.id = ct.tag_id
		where ct.course_id = $1
		`, v.ID)
		if err != nil {
			return nil, err
		}
		(*courses)[i].Tags = t
	}

	res := *courses
	return res, nil
}

func (s *CourseStore) Delete(ID int) error {
	_, err := s.db.Exec("delete from course where id =$1", ID)

	if err != nil {
		return err
	}
	return nil
}
