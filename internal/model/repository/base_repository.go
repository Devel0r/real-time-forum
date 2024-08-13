package repository

import "github.com/Pruel/real-time-forum/pkg/sqlite"

// MainRepository
type Repository struct {
    AuthRepo *AuthRepository
}

// New
func New(db *sqlite.Database) *Repository {
	return &Repository{
		AuthRepo: NewAuthRepository(db),
	}
}
