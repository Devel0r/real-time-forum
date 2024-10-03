package repository

import (
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

	res, err := c.DB.SQLite.Exec("INSERT INTO comment(id, content, user_id, post_id, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?)",
		comment.Id, comment.Content, comment.UserId, comment.PostId, comment.PostId, comment.CreatedAt, comment.UpdatedAt)
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

func (c *CommentRepository) GetCommentOwnerID(commentID int) (int, error) {
	var ownerID int
	err := c.DB.SQLite.QueryRow("SELECT user_id FROM comments WHERE id = ?", commentID).Scan(&ownerID)
	if err != nil {
		return 0, err
	}

	return ownerID, nil
}

func (c *CommentRepository) GetPostByID(postID int) (id int, err error) {

	err = c.DB.SQLite.QueryRow("SELECT id FROM post WHERE id=?").Scan(&postID)
	if err != nil {
		return 0, err
	}

	return postID, err
}

func (c *CommentRepository) DeleteComment(commentID int) error {
	if commentID == 0 {
		return serror.ErrEmptyCommentData
	}

	res, err := c.DB.SQLite.Exec("DELETE FROM comments WHERE id=?", commentID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return err
	}

	return nil
}
