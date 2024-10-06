package repository

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/Pruel/real-time-forum/internal/model"
	"github.com/Pruel/real-time-forum/pkg/serror"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
)

type CommentRepository struct {
	DB *sqlite.Database
}

func NewCommentRepository(db *sqlite.Database) *CommentRepository {
	return &CommentRepository{
		DB: db,
	}
}

func (c *CommentRepository) SaveComment(comment *model.Comment) (id int, err error) {
	if comment == nil {
		return 0, serror.ErrEmptyCommentData
	}

	res, err := c.DB.SQLite.Exec("INSERT INTO comments(content, user_id, post_id, created_at, updated_at) VALUES(?, ?, ?, ?, ?)",
		comment.Content, comment.UserId, comment.PostId, comment.PostId, comment.CreatedAt, comment.UpdatedAt)
	if err != nil {
		return 0, err
	}

	ID, err := res.LastInsertId()
	return int(ID), err

}

func (c *CommentRepository) GetUserSessionById(r *http.Request) (id int, err error) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return 0, nil
	}

	sessionUUID := cookie.Value

	var userID int
	err = c.DB.SQLite.QueryRow("SELECT user_id FROM session WHERE id=?", sessionUUID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, err
}

func (c *CommentRepository) GetCommentByCommentID(commentID int) (model.Comment, error) {
	com := model.Comment{}
	err := c.DB.SQLite.QueryRow("SELECT id, content, user_id, post_id, created_at, updated_at FROM comments WHERE id=?", commentID).Scan(&com.Id,
		&com.Content, &com.UserId, &com.PostId, &com.CreatedAt, &com.UpdatedAt)
	if err != nil {
		return com, err
	}

	return com, nil
}

func (c *CommentRepository) GetPostByID(postID int) (id int, err error) {

	err = c.DB.SQLite.QueryRow("SELECT id FROM post WHERE id=?").Scan(&postID)
	if err != nil {
		return 0, err
	}

	return postID, err
}

func (c *CommentRepository) DeleteComment(commentID int) error {
	_, err := c.DB.SQLite.Exec("DELETE FROM comments WHERE id=?", commentID)
	if err != nil {
		return err
	}

	return nil
}

// GetAllCommentsByPostID
func (c *CommentRepository) GetAllCommentsByPostID(postID int) (*[]model.Comment, error) {
	comments := &[]model.Comment{}
	crow, err := c.DB.SQLite.Query("SELECT id, content, user_id, post_id, created_at, updated_at FROM comments")
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Warn(err.Error())
			return comments, nil
		}
		slog.Error(err.Error())
		return nil, err
	}

	for crow.Next() {
		cm := model.Comment{}
		if err := crow.Scan(&cm.Id, &cm.Content, &cm.UserId, &cm.PostId, &cm.CreatedAt, &cm.UpdatedAt); err != nil {
			slog.Warn(err.Error())
			continue
		}

		*comments = append(*comments, cm)
	}

	return comments, nil
}
