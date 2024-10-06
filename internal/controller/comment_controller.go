package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Pruel/real-time-forum/internal/model"
	"github.com/Pruel/real-time-forum/internal/model/repository"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
)

// Comment Controller
type CommentController struct {
	CommRepo *repository.CommentRepository
}

func NewCommentController(db *sqlite.Database) *CommentController {
	return &CommentController{
		CommRepo: repository.NewCommentRepository(db),
	}
}

func (c *Controller) CreateComment(w http.ResponseWriter, r *http.Request) {
	// 1. Recieve data from front-end
	content := r.FormValue("comment-content")
	// 2. Validation for content
	if err := ValidateCommentContent(content); err != nil {
		slog.Warn("empty content field")
		ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	// must have: post_id, user_id
	// Надо получить юзера который написал комент и передать его в инстанс GetUserById
	userID, err := c.AuthController.ARepo.GetUserIDFromSession(w, r)
	if err != nil {
		slog.Error(err.Error())
		slog.Warn("Failed to get user by ID")
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}
	// И надо получить postID поста что бы так-же его передать в инстанс GetPostById
	postIDStr := r.FormValue("comment-post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		slog.Debug("Invalid post")
		ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	// 3. Create instanse of comment, and save external data to this instance
	currentTime := time.Now()
	comment := &model.Comment{
		Content:   content,
		UserId:    userID,
		PostId:    postID,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	// 4. Save into to database
	commentID, err := c.CommRepo.SaveComment(comment)
	if err != nil {
		slog.Error(err.Error())
		slog.Warn("Error when saving comments")
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	// 5. Log the change and redirect user to the same page
	fmt.Printf("Comment with ID: %d sucusessful create \n\n", commentID)

	rURL := fmt.Sprintf("/posts/%d", postID)
	fmt.Println("Redirect URL: ", rURL)
	http.Redirect(w, r, rURL, http.StatusSeeOther)
	return
}

func ValidateCommentContent(content string) error {
	if content == "" {
		return fmt.Errorf("Empty content")
	}

	return nil
}

func (c *Controller) DeleteComment(w http.ResponseWriter, r *http.Request) {
	// 1. Extract comment_id from QUERY
	commentIDStr := r.URL.Query().Get("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		slog.Debug("Invalid comment")
		ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	// 2. Extract userID for check ID for check comment owner with help if else {}
	userID, err := c.AuthController.ARepo.GetUserIDFromSession(w, r)
	if err != nil {
		slog.Error(err.Error())
		slog.Warn("Failed to get user by ID")
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	com, err := c.CommRepo.GetCommentByCommentID(commentID)
	if err != nil {
		slog.Error("Comment not found")
		ErrorController(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if com.UserId != userID {
		slog.Info("It is not your comment bratan")
		ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// 3. Remove comment
	if err := c.CommRepo.DeleteComment(commentID); err != nil {
		slog.Error("Failed to delete post")
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	fmt.Printf("Comment with ID: %d, deleted by user:%d\n", commentID, userID)

	rURL := fmt.Sprintf("/posts/%d", com.PostId)
	http.Redirect(w, r, rURL, http.StatusSeeOther)
	return
}
