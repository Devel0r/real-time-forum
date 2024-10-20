-- SQLBook: Code
BEGIN TRANSACTION;

CREATE TABLE users ( 
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    login TEXT NOT NULL UNIQUE, 
    age INTEGER NOT NULL, 
    gender TEXT NOT NULL DEFAULT 'Male',
    name TEXT NOT NULL,
    surname TEXT NOT NULL, 
    avatar TEXT,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL UNIQUE,
    rooms_id TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
    id TEXT NOT NULL, 
    user_id INTEGER NOT NULL REFERENCES users(id), 
    expires_at DATETIME NOT NULL DEFAULT 'now',
    created_at DATETIME NOT NULL DEFAULT 'now' 
);

CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT 'now'
);

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL, 
    image TEXT, 
    created_at DATETIME NOT NULL DEFAULT 'now', 
    updated_at DATETIME NOT NULL DEFAULT 'now',
    category_id INTEGER NOT NULL REFERENCES categories(id),
    user_id INTEGER NOT NULL REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    post_id INTEGER NOT NULL REFERENCES posts(id), 
    created_at DATETIME NOT NULL DEFAULT 'now', 
    updated_at DATETIME NOT NULL DEFAULT 'now'
);

CREATE TABLE IF NOT EXISTS rooms (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    client_creator_id INTEGER NOT NULL REFERENCES users(id),
    client_invited_id INTEGER NOT NULL REFERENCES users(id),
    last_msg TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS clients (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    username VARCHAR(255) NOT NULL,
    avatar TEXT,
    rooms_id TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    username VARCHAR(255) NOT NULL REFERENCES users(login),
    room_id VARCHAR(255) NOT NULL REFERENCES rooms(id),
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT 'now',
    is_delivered BOOLEAN,
    is_read BOOLEAN
);

INSERT INTO categories(title) VALUES('Anime'), ('Game'), ('Business'), ('Sport'), ('Life');

COMMIT; 