package store

import "github.com/jmoiron/sqlx"

type VideoStore struct {
	db *sqlx.DB
}

// TODO Haven't finished

func newVideoStore(DB *sqlx.DB) *VideoStore {
	return &VideoStore{
		db: DB,
	}
}

type UpdateVideoBody struct {
	FileName string
	ChapNum  int
}

type CreateVideoBody struct {
	FileName    string
	Description string
	UpdatedBy   int
}

// TODO : developing
func (s *VideoStore) CreateVideo(b CreateVideoBody) (int, error) {
	var lastInsertId int
	err := s.db.QueryRow("insert into video (file_name,description,updated_by) values ($1,$2,$3) returning id", b.FileName, b.Description, b.UpdatedBy).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (s *VideoStore) UpdateVideo(videoID int, b UpdateVideoBody) error {

	_, err := s.db.Exec("update video set file_name=$1 where id=$2", b.FileName, videoID)
	if err != nil {
		return err
	}
	return nil
}
