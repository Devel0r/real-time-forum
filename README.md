# Real Time 

- Фронтенд
- Бэкенд 
- Чат 

## SPA (Front-end)

## Forum Service

    1. Environment
        - configurations
        - linter
        - logger
        - autobuilding
        - CI (continue integration)
    2. Base logic
        - Main logic
        - App logic
    3. Layers
        - Transport (http server, router, handlers)
        - service (bisnes logic, main logic of the service - forum)
        - domen, core ( objeects, stuctures)
        - repository (main repository - conn and disconnection to the sqlite, CRUD AND OTHER DATA MANIPULATIONS)
    4. Features
        a. auth, User (Admin, Moderator, client)
            - identification (registration -> Nickname, Age, Gender, Name, Surname, Email, Password ) Sign up
                - validation, filtration
            - authontification - Sign in ( The user must be able to connect using either a nickname or e-mail combined with a password )
            - authorization (administration)
            - logout - Sign out ( The user should be able to log out from any page of the forum )
        b. Post ( date and time )
            - Create 
            - Update
            - Read
            - Delete
        c. Comment, (CRUD - basic operations, specefic operations )
            - must be parent post
            - create
            - update
            - read
            - delete
        d. Additional packages (logics)
            - sentinel errors
            - filtration
            - validation
            - encryting and decrypting
            - UUID

            MVC
                - Model ( Object, Struct - go, model)
                - View ( template, css, other )
                - Controller ()

        

## Chat Service
     Personal messages 

     a. Section to show who is online/offline and able to talk to

        - Organized by last message sent or alphabetically
        - Send private messages to online users
        - Section must be visible at all times

    b. Chat Section

        - Reload past messages when clicking on user
        - Show previous messages
        - Load last 10 messages, load more on scroll up without spamming

    c. Message Format

        - Date of the message sent
        - Username of the sender 
        

    d. Real-time functionality

        - Notify of new message without refreshing the page
        - Use WebSockets in backend and frontend 


<!-- Russian -->
## Сервис чата
     Личные сообщения

     a. Раздел для отображения тех, кто онлайн/в сети и с кем можно поговорить

        - Последовательность по последнему отправленному сообщению или в алфавитном порядке
        - Отправка личных сообщений онлайн-пользователям
        - Раздел должен быть постоянно видимым

    b. Раздел чата

        - Перезагрузка прошлых сообщений при нажатии на пользователя
        - Отображение предыдущих сообщений
        - Загрузка последних 10 сообщений, загрузка дополнительных при прокрутке вверх без спама

    c. Формат сообщений

        - Дата отправки сообщения
        - Имя пользователя отправителя

    d. Функциональность в реальном времени

        - Уведомление о новом сообщении без обновления страницы
        - Используйте WebSockets на бэкенде и фронтенде

### Chat service architerture, domains, features

#### DDD + Clean Arch maybe? //Go to white a house 


#### Domains

***Chat | Room*** 
    - id
    - chat_name
    - is_group_chat bool  (private, group)
    - **Clients** (userA, userB)
    - register 
    - unregister
    - **Message**
    - broadcast
    - **ChatRepository**

***Message***
    - id
    - user_or_sender
    - is_delivered
    - is_read 
    - content
    - date
    - reacts []*React - 
    - count_of_reacts -[CountEntity]Reacts

***Client***
    - id
    - username
    - other

    -- ws connnection
    -- chats
    -- messages


### Feature

 - create private chat
 - join chat by invite

 - get all chats by user id
 - get all online and all offline users of the list

 - write message
 - read message
 - reload message
 - get more messages ( infinity loads )

 Ч