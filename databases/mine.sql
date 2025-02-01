-- Create Users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT,
    password TEXT,
    email TEXT
);

-- Create Sessions table
CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    expires_at DATETIME
);

-- Create Posts table
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT,
    content TEXT,
    topic TEXT,
    like INTEGER DEFAULT 0,
    dislike INTEGER DEFAULT 0,
    commentcount INTEGER DEFAULT 0,
    create_at DATETIME NOT NULL
);

-- Create Comments table
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT,
    content TEXT,
    like INTEGER DEFAULT 0,
    dislike INTEGER DEFAULT 0,
    post_id INTEGER,
    create_at DATETIME,
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

-- Create Post Reactions table
CREATE TABLE IF NOT EXISTS post_reaction (
    score INTEGER DEFAULT 0,
    username TEXT,
    post_id INTEGER,
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

-- Create Comment Reactions table
CREATE TABLE IF NOT EXISTS comment_reaction (
    score INTEGER DEFAULT 0,
    username TEXT,
    comment_id INTEGER,
    FOREIGN KEY (comment_id) REFERENCES comments(id)
);
