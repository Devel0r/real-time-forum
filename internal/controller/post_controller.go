package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/Pruel/real-time-forum/internal/model"
	"github.com/Pruel/real-time-forum/internal/model/repository"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
)

type TData struct {
	Username   string
	Categories []model.Category
}

type PostController struct {
	message     string
	isValidData bool
	PostRepo    *repository.PostRepository
}

func NewPostController(db *sqlite.Database) *PostController {
	return &PostController{
		PostRepo: repository.NewPostRepository(db),
	}
}

// PostPage
func (m *Controller) CreatePage(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetTmpPath("post")))

	categories, err := m.PostController.PostRepo.GetAllCategories()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	userID, err := m.AuthController.ARepo.GetUserIDFromSession(w, r)
	if err != nil {
		slog.Warn(err.Error())
		http.Redirect(w, r, "/", http.StatusForbidden)
	}

	user, err := m.AuthController.ARepo.GetUserByUserID(userID)
	if err != nil {
		slog.Warn(err.Error())
		return
	}

	data := struct {
		User       *model.User
		Categories *[]model.Category
	}{
		User:       user,
		Categories: &categories,
	}

	if err := tmp.Execute(w, data); err != nil {
		slog.Error(err.Error())
		return
	}
}

// func CreatePost, HTTP Method POST -> Create
func (p *PostController) Create(w http.ResponseWriter, r *http.Request) {
	p.isValidData = true

	// 1. recieve data from front-end
	title := r.FormValue("post-title")
	content := r.FormValue("post-content")
	category := r.FormValue("post-category")

	// 2. validation
	if err := ValidatePostData(title, content); err != nil {
		slog.Warn("empty content field")
		p.isValidData = false
	}

	// 3. requirments
	userId, err := p.PostRepo.GetUserIdFromSession(r)
	if err != nil {
		slog.Error(err.Error())
		slog.Info("Failed to get user id from session")
		p.isValidData = false
		http.Redirect(w, r, "/sign-in", http.StatusUnauthorized)
		return
	}

	if !p.isValidData {
		slog.Warn("redirected user to create page")
		http.Redirect(w, r, "/create-posts", http.StatusBadRequest)
		return
	}

	// 4. create a instance of post, b) and save external data to this instance
	currentTime := time.Now()
	categoryId, _ := strconv.Atoi(category)
	post := &model.Post{
		Title:      title,
		Content:    content,
		CreatedAt:  currentTime,
		UpdatedAt:  currentTime,
		CategoryId: categoryId,
		UserId:     userId,
	}

	// 5. save into database
	postId, err := p.PostRepo.SavePost(post)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	// 6. redirect to main page
	fmt.Printf("Post with ID: %d sucusessful create \n\n", postId)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (p *PostController) Delete(w http.ResponseWriter, r *http.Request) {
	// Извлечь пост из запроса
	postIdStr := r.URL.Query().Get("post_id")
	postId, err := strconv.Atoi(postIdStr)
	if err != nil || postId < 0 {
		slog.Debug("Invalid post")
		return
	}

	// Извлечь юзера по ID
	userID, err := p.PostRepo.GetUserIdFromSession(r)
	if err != nil {
		slog.Error(err.Error())
		slog.Info("Failed to get user id from session")
		return
	}

	// Извлекаем пост из БД
	postID, err := p.PostRepo.GetPostByID(postId)
	if err != nil {
		slog.Error(err.Error())
		slog.Info("Post not found")
		return
	}

	// Проверка на то что этот юзер действительно владеет постом
	if postID != userID {
		slog.Error(err.Error())
		slog.Info("It is not your post bratan")
		return
	}

	// Удаляем пост
	if err = p.PostRepo.DeletePost(postId); err != nil {
		slog.Error("Failed to delete post")
		return
	}

	slog.Info("Post with ID: %d, deleted by user:%d\n")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ValidatePostData(title, content string) error {
	if title == "" && content == "" {
		return fmt.Errorf("The title and content of the post cannot be empty")
	}
	if title == "" {
		return fmt.Errorf("post title cannot be empty")
	}
	if content == "" {
		return fmt.Errorf("post content cannot be empty")
	}
	return nil
}

// mehtod get
// GetPostByID
