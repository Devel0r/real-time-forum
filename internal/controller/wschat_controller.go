package controller

import (
	"html/template"
	"log/slog"
	"net/http"
)

// WSChatController
type WsChatController struct {
}

// NewWSChatController constructor
func NewWSChatController() *WsChatController {
	return &WsChatController{}
}

// ChatPage
func (ws *Controller) ChatPage(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetTmpPath("wschat")))

	// get user
	userID, err := ws.AuthController.ARepo.GetUserIDFromSession(w, r)
	if err != nil {
		slog.Warn(err.Error())
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	user, err := ws.AuthController.ARepo.GetUserByUserID(userID)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	data := struct {
		Username string
	}{
		Username: user.Login,
	}

	if err := tmp.Execute(w, data); err != nil {
		slog.Error(err.Error())
		return
	}
}

// TODO: Group Chats
// CreateRoom | CreateGroupChat

// JoinRoom | JoinGroupChat

// GetRooms | GetGroupChats

// TODO: Private Chats
// CreatePvChat

// JoinPvChat

// GetPvChats

// TODO: Additional feature
