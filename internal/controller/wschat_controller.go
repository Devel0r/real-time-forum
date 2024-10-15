package controller

import (
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	wschat "github.com/Pruel/real-time-forum/internal/controller/ws_chat"
	"github.com/Pruel/real-time-forum/internal/model/repository"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
	"github.com/gorilla/websocket"
)

// WSChatController
type WsChatController struct {
	Hub   *wschat.ChatHub
	ARepo *repository.AuthRepository
	ChatRepo *repository.ChatRepository
}

// NewWSChatController constructor
func NewWSChatController(db *sqlite.Database) *WsChatController {
	return &WsChatController{
		Hub:   wschat.NewChat(),
		ARepo: repository.NewAuthRepository(db),
		ChatRepo: repository.NewChatReposotory(db),
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
	komnata := &wschat.Room{
		ID:      wschat.GetUUID(),
		Name:    room,
		Clients: make(map[string]*wschat.Client),
	}

	// user
	// pvchat | userA (creator) + userB, invite
	bibaID := strconv.Itoa(user.Id)
	Biba := wschat.Client{
		ID:       bibaID,
		Username: user.Login,
		RoomID:   komnata.ID,
		Conn:     conn,
		Message:  make(chan *wschat.Message, 10),
	}

	bobaID := strconv.Itoa(inUser.Id)
	Boba := wschat.Client{
		ID:       bobaID,
		Username: inUser.Login,
		RoomID:   komnata.ID,
		Conn:     conn,
		Message:  make(chan *wschat.Message, 10),
	}

	komnata.Clients[Biba.ID] = &Biba
	komnata.Clients[Boba.ID] = &Boba

	// add this room into rooms of the chat hub
	ws.Hub.Rooms[komnata.ID] = komnata

	// save this room into database
	room, err = ws.ChatRepo.SavePvChat()
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	// send this room to user, json
	aroom := wschat.SRoom{
		ID:      komnata.ID,
		Name:    komnata.Name,
		Clients: make(map[string]wschat.SClient),
	}

	aroom.Clients[bibaID] = wschat.SClient{
		ID:       bibaID,
		Username: Biba.Username,
		RoomID:   komnata.ID,
	}

	aroom.Clients[bobaID] = wschat.SClient{
		ID:       bobaID,
		Username: Boba.Username,
		RoomID:   komnata.ID,
	}

	if err := conn.WriteJSON(aroom); err != nil {
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
