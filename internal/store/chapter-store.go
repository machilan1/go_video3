package store

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type ChapterStore struct {
	db *sqlx.DB
}

type UpdateChapterBody struct {
	Title       string
	Description string
	ChapNum     int
}

type CreateChapterBody struct {
	Title       string
	Description string
	ChapNum     int
	VideoID     int
}

type Chapter struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	ChapNum     int       `db:"chap_num"`
	Length      int       `db:"length"`
	VideoID     int       `db:"video_id"`
	CourseID    int       `db:"course_id"`
	FileName    string    `db:"file_name"`
}

func newChapterStore(DB *sqlx.DB) *ChapterStore {
	return &ChapterStore{
		DB,
	}
}

// TODO : developing

func (s *ChapterStore) Create(courseID int, b CreateChapterBody) (int, error) {
	var chapterId int
	err := s.db.QueryRow("insert into chapter (title,description,chap_num,course_id,video_id) values($1,$2,$3,$4,$5) returning id", b.Title, b.Description, b.ChapNum, courseID, b.VideoID).Scan(&chapterId)
	if err != nil {
		return 0, err
	}
	return chapterId, nil
}

func (s *ChapterStore) FindOneChapter(chapterID int) (Chapter, error) {

	var chapter Chapter

	err := s.db.Get(&chapter, `
	select
		c.id,
		c.title,
		c.description,
		c.created_at,
		c.course_id,
		c.video_id,
		c.chap_num,
		v.file_name
	from chapter c
	left join video v on v.id= c.video_id 
	where c.id = $1
	 `, chapterID)
	if err != nil {
		return chapter, err
	}

	return chapter, err
}

func (s *ChapterStore) UpdateChapterInfo(chapterID int, b UpdateChapterBody) error {

	interpolationIndex := 1
	fieldListTemp := []string{}
	fieldList := []string{}
	valueList := []any{}

	if b.Description != "" {
		fieldListTemp = append(fieldListTemp, "description=")
		valueList = append(valueList, b.Description)
	}

	if b.Title != "" {
		fieldListTemp = append(fieldListTemp, "title=")
		valueList = append(valueList, b.Title)
	}

	if b.ChapNum != 0 {
		fieldListTemp = append(fieldListTemp, "chap_num=")
		valueList = append(valueList, b.ChapNum)
	}

	if len(fieldListTemp) < 1 {
		return fmt.Errorf("no field to update")
	}

	for _, v := range fieldListTemp {
		index := strconv.Itoa(interpolationIndex)
		fieldList = append(fieldList, v+`$`+index)
		interpolationIndex++
	}

	fieldsString := strings.Join(fieldList, ",")

	valueList = append(valueList, chapterID)
	index := strconv.Itoa(interpolationIndex)
	_, err := s.db.Exec(`update chapter set `+fieldsString+` where id =$`+index, valueList...)
	if err != nil {
		return err
	}
	return nil
}

func (s *ChapterStore) Delete(chapterID int) error {
	_, err := s.db.Exec("delete from chapter where id = $1", chapterID)
	if err != nil {
		return err
	}
	return nil
}
