package repository

import "github.com/Pruel/real-time-forum/pkg/sqlite"

type PostRepository struct {
	DB *sqlite.Database
}

func NewPostRepository(db *sqlite.Database) *PostRepository {
	return &PostRepository{
		DB: db,
	}
}
