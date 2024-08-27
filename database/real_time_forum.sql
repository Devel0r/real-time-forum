BEGIN TRANSACTION; 

-- fix(database/sql insts): change timestamp type on datetime, and current_timestamp on 'now"
--
--
--
CREATE TABLE IF NOT EXISTS 'users' ( 
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    login TEXT NOT NULL, 
    age INTEGER NOT NULL, 
    gender TEXT NOT NULL DEFAULT,
    name TEXT NOT NULL DEFAULT,
    surname TEXT NOT NULL DEFAULT, 
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id), 
    expired_at DATETIME NOT NULL DEFAULT 'now',
    created_at DATETIME NOT NULL DEFAULT 'now' 
);

CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL UNIQUE,
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




-- What needed to do: 
    -- private_messages


COMMIT; 