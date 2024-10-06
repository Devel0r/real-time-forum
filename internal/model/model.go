package model

import (
	"database/sql"
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
	Id        string    `sql:"id"`
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
	Id         int            `sql:"id"`
	Title      string         `sql:"title"`
	Content    string         `sql:"content"`
	Image      sql.NullString `sql:"image"`
	CreatedAt  time.Time      `sql:"created_at"`
	UpdatedAt  time.Time      `sql:"updated_at"`
	CategoryId int            `sql:"category_id"` // page+categories list  -> games -> games == db - games - id - 3 -> post.CategoryId = 3
	UserId     int            `sql:"user_id"`
}

// Если нам нужны какие либо манипулияции со временем, то юзаем time.Duration, если тупо временная метка то юзаемт time.Time
// time.Time -> 10:10:10s,

// time.Duration -> 10 second, i hour httpServer.WriteTimeout - 10 second

type Comment struct {
	Id        int    `sql:"id"`
	Content   string `sql:"content"`
	Author    string
	UserId    int       `sql:"user_id"`
	PostId    int       `sql:"post_id"`
	CreatedAt time.Time `sql:"created_at"`
	UpdatedAt time.Time `sql:"updated_at"`
}
