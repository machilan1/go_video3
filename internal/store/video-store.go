package store

import "github.com/jmoiron/sqlx"

type VideoStore struct {
	*sqlx.DB
}

func newVideoStore(DB *sqlx.DB) *VideoStore {
	return &VideoStore{
		DB,
	}
}

type UpdateVideoBody struct {
	FileName string
	ChapNum  int
}

func (s *VideoStore) UpdateVideo(videoID int, b UpdateVideoBody) error {

	_, err := s.Exec("update video set file_name = $1  where id =$2", b.FileName, videoID)
	if err != nil {
		return err
	}

	return nil

}
