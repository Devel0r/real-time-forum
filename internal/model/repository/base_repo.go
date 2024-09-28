package repository

import "github.com/Pruel/real-time-forum/pkg/sqlite"

// MainRepository
type Repository struct {
	AuthRepo *AuthRepository
	PostRepo *PostRepository
}

// New
func New(db *sqlite.Database) *Repository {
	return &Repository{
		AuthRepo: NewAuthRepository(db),
	}
}
