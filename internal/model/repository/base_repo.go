package repository

import (
	"github.com/Pruel/real-time-forum/pkg/sqlite"
)

// MainRepository
type Repository struct {
	AuthRepo    *AuthRepository
	PostRepo    *PostRepository
	CommentRepo *CommentRepository
	ChatRepo    *ChatRepository
}

// New
func New(db *sqlite.Database) *Repository {
	return &Repository{
		AuthRepo:    NewAuthRepository(db),
		PostRepo:    NewPostRepository(db),
		CommentRepo: NewCommentRepository(db),
		ChatRepo:    NewChatReposotory(db),
	}
}
