package model

import (
	"time"
)

type User struct {
	Id           int    `sql:"id"`
	Login        string `sql:"login"`
	Age          int    `sql:"age"`
	Gender       string `sql:"gender"`
	Name         string `sql:"name"`
	Surname      string `sql:"surname"`
	Email        string `sql:"email"`
	PasswordHash string `sql:"password_hash"`
}

type Session struct {
	Id        int       `sql:"id"`
	UserId    int       `sql:"user_id"`
	ExpiredAt time.Time `sql:"expired_at"`
	CreatedAt time.Time `sql:"created_at"`
}

type Category struct {
	Id        int       `sql:"id"`
	Title     string    `sql:"title"`
	CreatedAt time.Time `sql:"created_at"`
}

type Post struct {
	Id         int       `sql:"id"`
	Title      string    `sql:"title"`
	Content    string    `sql:"content"`
	Image      string    `sql:"image"`
	CreatedAt  time.Time `sql:"created_at"`
	UpdatedAt  time.Time `sql:"updated_at"`
	CategoryId int       `sql:"category_id"`
	UserId     int       `sql:"user_id"`
}

type Comment struct {
	Id        int       `sql:"id"`
	Content   string    `sql:"content"`
	UserId    int       `sql:"user_id"`
	PostId    int       `sql:"post_id"`
	CreatedAt time.Time `sql:"created_at"`
	UpdatedAt time.Time `sql:"updated_at"`
}
