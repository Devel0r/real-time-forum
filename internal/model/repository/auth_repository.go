package repository

import "github.com/Pruel/real-time-forum/pkg/sqlite"

type AuthRepository struct {
    DB *sqlite.Database

}

func NewAuthRepository(db *sqlite.Database) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

// SignUp
    // create session

// SignIn
    // check session

// SignOut
    // remove session