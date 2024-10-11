package controller

import (
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Pruel/real-time-forum/internal/model/repository"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
	"github.com/gorilla/websocket"
)

// WSChatController
type WsChatController struct {
	Hub   *ChatHub
	ARepo *repository.AuthRepository
	// WsChatRepo *repository.WsChatRepository
}

// NewWSChatController constructor
func NewWSChatController(db *sqlite.Database) *WsChatController {
	return &WsChatController{
		ARepo: repository.NewAuthRepository(db),
	}
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
func (ws *WsChatController) CreateRoom(w http.ResponseWriter, r *http.Request) {
	// create a new pv chat frontend -> event  js getElementByID -> onclick ->  method new websocket-> get ws://localhost:80/roomID/3/?usernameA=Janika&usernameB=Daniil
	upgrader := websocket.Upgrader{
		WriteBufferSize: 1024, // 1024 byte = 1 mb // Хорошая фича
		ReadBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			// origin := r.Header
			// remIP := r.RemoteAddr
			// url := r.URL

			// if url != "http://localhost:8080" {
			// 	return false
			// }
			return true
		},
	}

	// switch http to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	room := r.URL.Query().Get("room") // room name
	inviteUsername := r.URL.Query().Get("invited_username")
	idstr := r.URL.Query().Get("userID")
	userID, _ := strconv.Atoi(idstr)

	// ws://localhost:8081/ws/create-pvchat?room=amongASS&userID=1&invited_username=simpleTest

	// params
	// userID, err := ws.ARepo.GetUserIDFromSession(w, r)
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
	// 	return
	// }

	user, err := ws.ARepo.GetUserByUserID(userID)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	inUser, err := ws.ARepo.GetUserByUsername(inviteUsername)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	// create a new instance of the room
	komnata := &Room{
		ID:      GetUUID(),
		Name:    room,
		Clients: make(map[string]*Client),
	}

	// user
	// pvchat | userA (creator) + userB, invite
	bibaID := strconv.Itoa(user.Id)
	Biba := Client{
		ID:       bibaID,
		Username: user.Login,
		RoomID:   komnata.ID,
		Conn:     conn,
		Message:  make(chan *Message, 10),
	}

	bobaID := strconv.Itoa(inUser.Id)
	Boba := Client{
		ID:       bobaID,
		Username: inUser.Login,
		RoomID:   komnata.ID,
		Conn:     conn,
		Message:  make(chan *Message, 10),
	}

	komnata.Clients[Biba.ID] = &Biba
	komnata.Clients[Boba.ID] = &Boba

	// add this room into rooms of the chat hub
	ws.Hub.Rooms[komnata.ID] = komnata

	// save this room into database
	// TODO: save into database

	// send this room to user, json

	if err := conn.WriteJSON(komnata); err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

// JoinRoom | JoinGroupChat
func (c *WsChatController) JoinRoom() {

}

// GetRooms | GetGroupChats
func (c *WsChatController) GetRooms() {

}

// GetOnlineAndOfflineUsers

// TODO: Private Chats
// CreatePvChat

// JoinPvChat

// GetPvChats

// TODO: Additional feature
