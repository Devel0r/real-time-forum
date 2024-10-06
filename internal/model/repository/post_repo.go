package repository

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/Pruel/real-time-forum/internal/model"
	"github.com/Pruel/real-time-forum/pkg/serror"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
)

type PostRepository struct {
	DB *sqlite.Database
}

func NewPostRepository(db *sqlite.Database) *PostRepository {
	return &PostRepository{
		DB: db,
	}
}

func (p *PostRepository) SavePost(post *model.Post) (id int, err error) {
	if post == nil {
		return 0, serror.ErrEmptyPostData
	}

	res, err := p.DB.SQLite.Exec("INSERT INTO posts(title, content, created_at, updated_at, category_id, user_id) VALUES(?, ?, ?, ?, ?, ?)",
		post.Title, post.Content, post.CreatedAt, post.UpdatedAt, post.CategoryId, post.UserId)
	if err != nil {
		return 0, err
	}

	ID, _ := res.LastInsertId()
	return int(ID), err
}

func (p *PostRepository) DeletePost(postId int) error {
	if postId == 0 {
		return serror.ErrEmptyPostData
	}

	res, err := p.DB.SQLite.Exec("DELETE FROM posts WHERE id=?", postId)
	if err != nil {
		return err
	}

	// Проверяем, сколько строк было затронуто
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Пост не был удалён, возможно, он не найден
		return err
	}

	return nil
}

func (p *PostRepository) GetPostByID(postId int) (id int, err error) {

	err = p.DB.SQLite.QueryRow("SELECT id FROM post WHERE id=?").Scan(&postId)
	if err != nil {
		return 0, nil
	}

	return postId, nil
}

func (p *PostRepository) GetUserIdFromSession(r *http.Request) (int, error) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return 0, err
	}

	sessionUUID := cookie.Value

	var userId int
	err = p.DB.SQLite.QueryRow("SELECT user_id FROM sessions WHERE id=?", sessionUUID).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

// GetAllCategories
func (p *PostRepository) GetAllCategories() ([]model.Category, error) {
	crows, err := p.DB.SQLite.Query("SELECT id, title, created_at FROM categories")
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Warn("Categories not found")
		}
		slog.Error(err.Error())
		return nil, err
	}

	categories := []model.Category{}
	for crows.Next() {
		category := model.Category{}
		if err := crows.Scan(&category.Id, &category.Title, &category.CreatedAt); err != nil {
			if err == sql.ErrNoRows {
				slog.Warn("Category not found")
				continue
			}
			slog.Error(err.Error())
			break
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// GetCategoryByCategoryID
func (p *PostRepository) GetCategoryByCategoryID(id int) (*model.Category, error) {
	cat := model.Category{}
	if err := p.DB.SQLite.QueryRow("SELECT id, title, created_at FROM categories WHERE id=?", id).Scan(&cat.Id, &cat.Title, &cat.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			slog.Warn(err.Error())
			cat.Title = "Unknown"
		}
	}

	return &cat, nil
}

// GetAllPosts
func (p *PostRepository) GetAllPosts() (*[]model.Post, error) {
	prow, err := p.DB.SQLite.Query("SELECT id, title, content, image, created_at, updated_at, category_id, user_id FROM posts")
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Warn(err.Error())
		}
		slog.Error(err.Error())
		return nil, err
	}

	posts := []model.Post{}
	for prow.Next() {
		post := model.Post{}
		err := prow.Scan(&post.Id, &post.Title, &post.Content, &post.Image, &post.CreatedAt, &post.UpdatedAt, &post.CategoryId, &post.UserId)
		if err != nil {
			if err == sql.ErrNoRows {
				slog.Warn(err.Error())
				continue
			}
			slog.Error(err.Error())
			break
		}
		posts = append(posts, post)
	}

	return &posts, nil
}
