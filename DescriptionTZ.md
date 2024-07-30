    Objectives

On this project you will have to focus on a few points:

    Registration and Login
    Creation of posts
        Commenting posts
    Private Messages

As you already did the first forum you can use part of the code, but not all of it. Your new forum will have five different parts:

    SQLite, in which you will store data, just like in the previous forum
    Golang, in which you will handle data and Websockets (Backend)
    Javascript, in which you will handle all the Frontend events and clients Websockets
    HTML, in which you will organize the elements of the page
    CSS, in which you will stylize the elements of the page

You will have only one HTML file, so every change of page you want to do, should be handled in the Javascript. This can be called having a single page application.
Registration and Login

To be able to use the new and upgraded forum users will have to register and login, otherwise they will only see the registration or login page. This is premium stuff. The registration and login process should take in consideration the following features:

    Users must be able to fill a register form to register into the forum. They will have to provide at least:
        Nickname
        Age
        Gender
        First Name
        Last Name
        E-mail
        Password
    The user must be able to connect using either the nickname or the e-mail combined with the password.
    The user must be able to log out from any page on the forum.

Posts and Comments

This part is pretty similar to the first forum. Users must be able to:

    Create posts
        Posts will have categories as in the first forum
    Create comments on the posts
    See posts in a feed display
        See comments only if they click on a post

Private Messages

Users will be able to send private messages to each other, so you will need to create a chat, where it will exist :

    A section to show who is online/offline and able to talk to:
        This section must be organized by the last message sent (just like discord). If the user is new and does not present messages you must organize it in alphabetic order.
        The user must be able to send private messages to the users who are online.
        This section must be visible at all times.

    A section that when clicked on the user that you want to send a message, reloads the past messages. Chats between users must:
        Be visible, for this you will have to be able to see the previous messages that you had with the user
        Reload the last 10 messages and when scrolled up to see more messages you must provide the user with 10 more, without spamming the scroll event. Do not forget what you learned!! (Throttle, Debounce)

    Messages must have a specific format:
        A date that shows when the message was sent
        The user name, that identifies the user that sent the message

As it is expected, the messages should work in real time, in other words, if a user sends a message, the other user should receive the notification of the new message without refreshing the page. Again this is possible through the usage of WebSockets in backend and frontend.
Allowed Packages

    All standard go packages are allowed.
    Gorilla websocket
    sqlite3
    bcrypt
    UUID

    You must not use use any frontend libraries or frameworks like React, Angular, Vue etc.

This project will help you learn about:

    The basics of web :
        HTML
        HTTP
        Sessions and cookies
        CSS
        Backend and Frontend
        DOM
    Go routines
    Go channels
    WebSockets:
        Go Websockets
        JS Websockets
    SQL language
        Manipulation of databases


    Задачи

В этом проекте вам предстоит сосредоточиться на нескольких моментах:

    Регистрация и вход
    Создание постов
        Комментирование постов
    Личные сообщения

Поскольку вы уже сделали первый форум, вы можете использовать часть кода, но не весь. Ваш новый форум будет состоять из пяти различных частей:

    SQLite, в которой вы будете хранить данные, как и в предыдущем форуме
    Golang, в котором вы будете работать с данными и Websockets (Backend)
    Javascript, в котором вы будете обрабатывать все события Frontend и клиенты Websockets
    HTML, в котором вы будете организовывать элементы страницы
    CSS, в котором вы будете стилизовать элементы страницы

У вас будет только один HTML-файл, поэтому каждое изменение страницы, которое вы хотите сделать, должно быть обработано в Javascript. Это можно назвать одностраничным приложением.
Регистрация и вход

Чтобы иметь возможность пользоваться новым и обновленным форумом, пользователи должны будут зарегистрироваться и войти, иначе они будут видеть только страницу регистрации или входа. Это премиум-функции. Процесс регистрации и входа должен учитывать следующие особенности:

    Пользователи должны иметь возможность заполнить регистрационную форму для регистрации на форуме. Они должны будут указать как минимум:    
        Никнейм
        Возраст
        Пол
        Имя
        Фамилия
        E-mail
        Пароль
    Пользователь должен иметь возможность подключаться, используя либо ник, либо e-mail в сочетании с паролем.
    Пользователь должен иметь возможность выйти из системы с любой страницы форума.

Сообщения и комментарии

Эта часть практически аналогична первому форуму. Пользователи должны иметь возможность:

    Создавать сообщения
        У сообщений будут категории, как и в первом форуме
    Создавать комментарии к сообщениям
    Просматривать сообщения в виде ленты
        Видеть комментарии только при нажатии на сообщение

Личные сообщения

Пользователи смогут отправлять личные сообщения друг другу, поэтому вам нужно будет создать чат, где он будет существовать:

    Раздел, показывающий, кто находится онлайн/оффлайн и с кем можно поговорить:
        Этот раздел должен быть организован по последнему отправленному сообщению (как в Discord). Если пользователь новый и у него нет сообщений, вы должны организовать его в алфавитном порядке.
        Пользователь должен иметь возможность отправлять личные сообщения пользователям, которые находятся онлайн.
        Этот раздел должен быть виден всегда.

    Раздел, который при нажатии на пользователя, которому вы хотите отправить сообщение, перезагружает прошлые сообщения. Чаты между пользователями должны:
        Быть видимыми, для этого вы должны иметь возможность видеть предыдущие сообщения, которые были у вас с пользователем.
        Перезагружать последние 10 сообщений и при прокрутке вверх, чтобы увидеть больше сообщений, вы должны предоставить пользователю еще 10, не спамя событие прокрутки. Не забывайте, чему вы научились!!! (Throttle, Debounce)

    Сообщения должны иметь определенный формат:
        Дата, которая показывает, когда сообщение было отправлено
        Имя пользователя, которое идентифицирует пользователя, отправившего сообщение.

Как и ожидалось, сообщения должны работать в реальном времени, другими словами, если пользователь отправляет сообщение, другой пользователь должен получить уведомление о новом сообщении, не обновляя страницу. Это возможно благодаря использованию WebSockets в бэкенде и фронтенде.
Разрешенные пакеты

    Разрешены все стандартные пакеты go.
    Gorilla websocket
    sqlite3
    bcrypt
    UUID

    Вы не должны использовать какие-либо библиотеки или фреймворки для фронтенда, такие как React, Angular, Vue и т.д.

Этот проект поможет вам узнать о:

    Основы веб-технологий:
        HTML
        HTTP
        Сессии и куки
        CSS
        Backend и Frontend
        DOM
    Процедуры Go
    Каналы Go
    WebSockets:
        Go Websockets
        JS Websockets
    Язык SQL
        Манипулирование базами данных