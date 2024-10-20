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
	ClientsList        []model.User // all online and offline users
	Messages           []wschat.Message // all messages by the current room
	CurrentRoomClients []wschat.SRoom   // all the current client rooms
	Username           string
}

// getClinets
func (ws *WsChatController) getClients(chat *wschat.ChatHub, dbClients []model.User) (cls []model.User, err error) {
	if chat == nil || dbClients == nil {
		return nil, errors.New("error, nil struct pointer")
	}

	for _, wcl := range chat.Clients {
		for _, dcl := range dbClients {
			if wcl.Username == dcl.Login {
				dcl.IsOnline = true
			}
			cls = append(cls, dcl)
		}
	}

	return cls, nil
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

	// TODO:  depricated,  remove this code
	// client, err := ws.ChatRepo.GetClientByUsername(user.Login)
	// if err != nil {
	// 	return chatData, err
	// }

	rooms, err := ws.ChatRepo.GetAllRoomsByClientID(user.Id)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return chatData, err
	}
	chatData.CurrentRoomClients = rooms

	// TODO: depricated, remove this code
	// clients, err := ws.ChatRepo.GetAllClients()
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	// 	return chatData, err
	// }

	users, err := ws.ARepo.GetAllUsers()
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return chatData, err
	}

	// TODO: implement last message saving for room.last_message

	cls, err := ws.getClients(ws.Hub, users)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return chatData, err
	}
	chatData.ClientsList = cls

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
	// params

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
	}

	Boba := wschat.Client{
		ID:       inUser.Id,
		Username: inUser.Login,
		RoomID:   komnata.ID,
		Conn:     conn, // TODO: maybe a bug
		Message:  make(chan wschat.Message, 10),
	}

	// add this room into rooms of the chat hub
	ws.Hub.Rooms[komnata.ID] = komnata

	// send this room to user, json
	aroom := wschat.SRoom{
		ID:      komnata.ID,
		Name:    komnata.Name,
	}

	// TODO: depricated, remove this code
	// biba := wschat.SClient{
	// 	ID:       bibaID,
	// 	Username: Biba.Username,
	// 	RoomID:   []string{komnata.ID},
	// }
	// aroom.Clients[bibaID] = biba

	// boba := wschat.SClient{
	// 	ID:       bobaID,
	// 	Username: Boba.Username,
	// 	RoomID:   []string{komnata.ID},
	// }
	// aroom.Clients[bobaID] = boba

	gMsg := wschat.Message{
		ID:      wschat.GetUUID(),
		From:    Biba.Username,
		RoomID:  komnata.ID,
		Content: fmt.Sprintf("the %s created the %s chat, and add the %s\n", Biba.Username, aroom.Name, Boba.Username),
	}
	aroom.ClientCretorID = user.Id
	aroom.ClientCretorID = inUser.Id
	aroom.LastMessage = &gMsg // TODO: meybe bug

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

	// TODO: depricated, reomve this code
	// clIDs, err := ws.ChatRepo.SaveClients(&biba, &boba)
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	// 	return
	// }
	// slog.Info("clients successful saved into dabase", "clients_ids", strings.Join(clIDs, ", "))

	if err := conn.WriteJSON(aroom); err != nil {
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
		Message: make(chan wschat.Message, 10),
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

	// TODO: depricated, remove this code
	// dclient, err := c.ChatRepo.GetClientByUsername(username)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		cl := wschat.SClient{
	// 			ID:       wschat.GetUUID(),
	// 			Username: username,
	// 			RoomID:   roomID,
	// 		}

	// 		Client.ID = cl.ID
	// 		Client.Username = cl.Username
	// 		Client.RoomID = cl.RoomID

	// 		id, err := c.ChatRepo.SaveClients(&cl)
	// 		if err != nil {
	// 			slog.Error(errors.New("error, wrong user in chat").Error(), "err", err.Error())
	// 			ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	// 			return
	// 		}
	// 		slog.Info("a new joining the chat client successful saved into db", "client_id", strings.Join(id, ", "))
	// 	}
	// 	slog.Error(err.Error())
	// 	ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	// 	return
	// }
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

	_, err = c.ChatRepo.SaveMessage(&gMsg)
	if err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if err := c.ChatRepo.SaveLasgRoomMsg(droom.ID, &gMsg); err != nil {
		slog.Error(err.Error())
		ErrorController(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	//9. send to front-end ws info about this operations
	Client.Conn.WriteJSON(gMsg)

	//6. broadcast this message in this room
	// TODO: test again Hub.Broadcast, and Register, with Unregister channels
	c.Hub.Broadcast <- gMsg
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
