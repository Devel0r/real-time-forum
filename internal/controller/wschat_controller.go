package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	wschat "github.com/Pruel/real-time-forum/internal/controller/wschat"
	"github.com/Pruel/real-time-forum/internal/model"
	"github.com/Pruel/real-time-forum/internal/model/repository"
	"github.com/Pruel/real-time-forum/pkg/serror"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
)

// WSChatController
type WsChatController struct {
	Hub      *wschat.ChatHub
	ARepo    *repository.AuthRepository
	ChatRepo *repository.ChatRepository
}

// NewWSChatController constructor
func NewWSChatController(db *sqlite.Database, ch *wschat.ChatHub) *WsChatController {
	return &WsChatController{
		Hub:      ch,
		ARepo:    repository.NewAuthRepository(db),
		ChatRepo: repository.NewChatReposotory(db),
	}
}

// ChatData
type ChatData struct {
	Users              []*model.User    // all online and offline users
	Messages           []wschat.Message // all messages by the current room
	CurrentRoomClients []wschat.SRoom   // all the current client rooms
	Username           string
	CurrentRoomID      string
}

var OnlineUsers = []*model.User{}

// getClinets
func (ws *WsChatController) getClients(chat *wschat.ChatHub, dbClients []*model.User, username string) ([]*model.User, error) {
	if chat == nil || dbClients == nil {
		return nil, errors.New("error, nil struct pointer")
	}

	for i, dcl := range dbClients {
		if dcl.Login == username {
			dcl.IsOnline = true
		}
		for _, wcl := range chat.Clients {
			if wcl.Username == dcl.Login {
				dbClients[i].IsOnline = true
			}
		}
	}

	return dbClients, nil
}

func (ws *WsChatController) getChatData(w http.ResponseWriter, r *http.Request) (ChatData, error) {
	chatData := ChatData{}
	if w == nil || r == nil {
		return chatData, errors.New("error, nil arguments")
	}

	userID, err := ws.ARepo.GetUserIDFromSession(w, r)
	if err != nil {
		slog.Warn(err.Error())
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return chatData, err
	}

	user, err := ws.ARepo.GetUserByUserID(userID)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return chatData, err
	}
	chatData.Username = user.Login

	rooms, err := ws.ChatRepo.GetAllRoomsByClientID(user.Id)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return chatData, err
	}
	chatData.CurrentRoomClients = rooms

	ChUsers, err := ws.ARepo.GetAllUsers()
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return chatData, err
	}

	cls, err := ws.getClients(ws.Hub, ChUsers, user.Login)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return chatData, err
	}
	chatData.Users = cls

	return chatData, nil
}

// ChatPage
func (ws *Controller) ChatPage(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetTmpPath("wschat")))

	// get chat data
	chatData, err := ws.WsChatController.getChatData(w, r)
	if err != nil {
		slog.Error(err.Error())
	}

	if err := tmp.Execute(w, chatData); err != nil {
		slog.Error(err.Error())
		return
	}
}

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
	// params

	fmt.Printf("\nRequest for create a new room, room: %s, invited_username: %s\n", room, inviteUsername)

	uID, err := ws.ARepo.GetUserIDFromSession(w, r)
	if err != nil {
		slog.Error(err.Error())
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	user, err := ws.ARepo.GetUserByUserID(uID)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	inUser, err := ws.ARepo.GetUserByUsername(inviteUsername)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, "Invited user don't exists")
		return
	}

	// create a new instance of the room
	komnata := &wschat.Room{
		ID:      wschat.GetUUID(),
		Name:    room,
		Clients: make(map[int]*wschat.Client),
	}

	// user
	// pvchat | userA (creator) + userB, invite
	Biba := wschat.Client{
		ID:       user.Id,
		Username: user.Login,
		RoomID:   komnata.ID,
		Conn:     conn,
		Message:  make(chan wschat.Message, 10),
		DB:       ws.ChatRepo.DB,
	}

	Boba := wschat.Client{
		ID:       inUser.Id,
		Username: inUser.Login,
		RoomID:   komnata.ID,
		Conn:     conn, // TODO: maybe a bug
		Message:  make(chan wschat.Message, 10),
		DB:       ws.ChatRepo.DB,
	}

	// add this room into rooms of the chat hub
	ws.Hub.Rooms[komnata.ID] = komnata

	// send this room to user, json
	aroom := wschat.SRoom{
		ID:   komnata.ID,
		Name: komnata.Name,
	}

	gMsg := wschat.Message{
		ID:      wschat.GetUUID(),
		From:    Biba.Username,
		RoomID:  komnata.ID,
		Content: fmt.Sprintf("the %s created the %s chat, and add the %s\n", Biba.Username, aroom.Name, Boba.Username),
	}
	aroom.ClientCretorID = user.Id
	aroom.ClientInvitedID = inUser.Id
	aroom.LastMessage = &gMsg

	msgID, err := ws.ChatRepo.SaveMessage(&gMsg)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	slog.Info("message successful saved into daatabase", "msgID", msgID)

	// save this room into database
	roomID, err := ws.ChatRepo.SaveRoom(&aroom)
	if err != nil {
		if !errors.Is(err, serror.ErrRoomAlreadyExists) {
			slog.Error(err.Error())
			ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}
	slog.Info("the room successful saved into database", "room_id", roomID)

	// information message for fron-end
	res := struct {
		Success  bool   `json:"success"`
		Room     string `json:"room"`
		RoomID   string `json:"roomID"`
		Username string `json:"username"`
		Message  string `json:"message"`
	}{
		Success:  true,
		Room:     aroom.Name,
		RoomID:   aroom.ID,
		Username: inUser.Login,
		Message:  gMsg.Content,
	}

	if err := conn.WriteJSON(res); err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	slog.Info(fmt.Sprintf("the %s chat successfull created by %s for %s, chat | room is active", komnata.Name, Biba.Username, Boba.Username))
	ws.Hub.Rooms[komnata.ID] = komnata
	ws.Hub.Register <- &Biba
	ws.Hub.Register <- &Boba
	ws.Hub.Broadcast <- gMsg

	go Biba.WriteMessage()
	Biba.ReadMessage(ws.Hub)
	slog.Warn("the %s chat is diactived by clients, websocket connection is broken")
}

// JoinRoom | JoinGroupChat
func (c *WsChatController) JoinRoom(w http.ResponseWriter, r *http.Request) {
	//1. create a new upgrader instance, make upgrade to ws
	upgrader := websocket.Upgrader{
		WriteBufferSize: 1024, // 1024 bytes = 1 mb
		ReadBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			slog.Error(err.Error())
			ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}

	// ws://localhost:8081/ws/join-room?room_id=?&username=daniil
	//2. recieve user data, after validating this data
	roomID := r.URL.Query().Get("room_id")
	username := r.URL.Query().Get("username")

	if roomID == "" || username == "" {
		slog.Error(errors.New("error, empty joining room data").Error())
		ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	//3. compare this room_id with room_id from database, ws.ChatHub also username, user_id
	droom, err := c.ChatRepo.GetRoomByID(roomID)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if droom.ID != roomID {
		slog.Error(errors.New("error, wrong room id").Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	userID, err := c.ARepo.GetUserIDFromSession(w, r)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user, err := c.ARepo.GetUserByUserID(userID)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if username != user.Login {
		slog.Error(errors.New("error, wrong user in chat").Error())
		ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	//4. if invited user not exists in database and chat service pool, register this user in chat pool
	// and save this user into database
	Client := wschat.Client{
		Conn:    conn,
		Message: make(chan wschat.Message, 10),
		DB:      c.ChatRepo.DB,
	}

	dclient, err := c.ARepo.GetUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Warn(err.Error())
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if droom.ClientInvitedID != dclient.Id && droom.ClientCretorID != dclient.Id {
		slog.Error(errors.New("error, uninvited client try to join the chat").Error())
		ErrorController(w, http.StatusMethodNotAllowed, "Access denied! Its not your chat!")
		return
	}
	slog.Info("the client joining the chat", "username", dclient.Login)

	if dclient.Login != "" {
		Client.ID = dclient.Id
		Client.Username = dclient.Login
		Client.RoomID = droom.ID
	}

	//5. create a new message for greating this user in chat
	gMsg := wschat.Message{
		ID:        wschat.GetUUID(),
		From:      Client.Username,
		RoomID:    roomID,
		Content:   fmt.Sprintf("the %s join the chat", Client.Username),
		CreatedAt: time.Now(),
	}

	if err := c.ChatRepo.SaveLasgRoomMsg(droom.ID, &gMsg); err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	// TODO: load all old this room messages by room_id
	messages, err := c.ChatRepo.GetAllMessagesByRoomID(droom.ID)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if len(messages) > 0 {
		for _, msg := range messages {
			Client.Conn.WriteJSON(msg)
		}
	}

	//9. send to front-end ws info about this operations
	// Client.Conn.WriteJSON(gMsg)

	//6. broadcast this message in this room
	// TODO: test again Hub.Broadcast, and Register, with Unregister channels
	room := wschat.Room{
		ID:      droom.ID,
		Name:    droom.Name,
		Clients: make(map[int]*wschat.Client),
	}
	room.Clients[Client.ID] = &Client
	c.Hub.Rooms[droom.ID] = &room
	// c.Hub.Broadcast <- gMsg
	c.Hub.Register <- &Client

	//7. add ws connection, client.conn

	//8. run client's write and read methods in separate goroutines
	go Client.WriteMessage()
	Client.ReadMessage(c.Hub)

	// js-> ws -> getRooms >all messages, users and other info
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
