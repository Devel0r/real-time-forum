package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	wschat "github.com/Pruel/real-time-forum/internal/controller/wschat"
	"github.com/Pruel/real-time-forum/internal/model/repository"
	"github.com/Pruel/real-time-forum/pkg/sqlite"
)

// WSChatController
type WsChatController struct {
	Hub      *wschat.ChatHub
	ARepo    *repository.AuthRepository
	ChatRepo *repository.ChatRepository
}

// NewWSChatController constructor
func NewWSChatController(db *sqlite.Database) *WsChatController {
	return &WsChatController{
		Hub:      wschat.NewChat(),
		ARepo:    repository.NewAuthRepository(db),
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

	// ws://localhost:8081/ws/create-pvchat?room=amongASS&userID=2&invited_username=user1

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

	// send this room to user, json
	aroom := wschat.SRoom{
		ID:      komnata.ID,
		Name:    komnata.Name,
		Clients: make(map[string]wschat.SClient),
	}

	biba := wschat.SClient{
		ID:       bibaID,
		Username: Biba.Username,
		RoomID:   komnata.ID,
	}
	aroom.Clients[bibaID] = biba

	boba := wschat.SClient{
		ID:       bobaID,
		Username: Boba.Username,
		RoomID:   komnata.ID,
	}
	aroom.Clients[bobaID] = boba

	// save this room into database
	roomID, err := ws.ChatRepo.SaveRoom(&aroom)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	slog.Info("the room successful saved into database", "room_id", roomID)

	clIDs, err := ws.ChatRepo.SaveClients(&biba, &boba)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	slog.Info("clients successful saved into dabase", "clients_ids", strings.Join(clIDs, ", "))

	if err := conn.WriteJSON(aroom); err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
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
	droom, err := c.ChatRepo.GetRoomID(roomID)
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

	// TODO: set session coockie

	// TODO: fix this
	// userID, err := c.ARepo.GetUserIDFromSession(w, r)
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	// 	return
	// }

	// user, err := c.ARepo.GetUserByUsername(username)
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	// 	return
	// }

	// // TODO: ref:
	// if username != username {
	// 	slog.Error(errors.New("error, wrong user in chat").Error())
	// 	ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	// 	return
	// }

	//4. if invited user not exists in database and chat service pool, register this user in chat pool
	// and save this user into database
	Client := wschat.Client{
		Conn:    conn,
		Message: make(chan *wschat.Message, 10),
	}

	dclient, err := c.ChatRepo.GetClientByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			cl := wschat.SClient{
				ID:       wschat.GetUUID(),
				Username: username,
				RoomID:   roomID,
			}

			Client.ID = cl.ID
			Client.Username = cl.Username
			Client.RoomID = cl.RoomID

			id, err := c.ChatRepo.SaveClients(&cl)
			if err != nil {
				slog.Error(errors.New("error, wrong user in chat").Error(), "err", err.Error())
				ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
				return
			}
			slog.Info("a new joining the chat client successful saved into db", "client_id", strings.Join(id, ", "))
		}
		slog.Error(err.Error())
		ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	slog.Info("the client joining the chat", "username", dclient.Username)

	if dclient.Username != "" {
		Client.ID = dclient.ID
		Client.Username = dclient.Username
		Client.RoomID = dclient.RoomID
	}

	//5. create a new message for greating this user in chat
	gMsg := wschat.Message{
		ID:        wschat.GetUUID(),
		From:      Client.Username,
		RoomID:    roomID,
		Content:   fmt.Sprintf("the %s join the chat", Client.Username),
		CreatedAt: time.Now(),
	}

	_, err = c.ChatRepo.SaveMessage(&gMsg)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	//9. send to front-end ws info about this operations
	Client.Conn.WriteJSON(gMsg)

	//6. broadcast this message in this room
	c.Hub.Broadcast <- &gMsg
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
