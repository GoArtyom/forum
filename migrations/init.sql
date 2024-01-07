CREATE TABLE IF NOT EXISTS users (
    id INTEGER NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    mode INTEGER NOT NULL,
    rols INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS sessions (
    user_id INTEGER NOT NULL,
    uuid TEXT NOT NULL,
    expire_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT,
    user_id INTEGER NOT NULL,
    user_name TEXT NOT NULL,
    create_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS posts_images(
    id INTEGER NOT NULL PRIMARY KEY,
    post_id INTEGER NOT NULL,
    name text NOT NULL,
    type text NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER NOT NULL PRIMARY KEY,
    post_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    user_name TEXT NOT NULL,
    create_at DATETIME NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS category (name TEXT NOT NULL UNIQUE PRIMARY KEY);

CREATE TABLE IF NOT EXISTS posts_categories (
    post_id INTEGER NOT NULL,
    category_name TEXT NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_name) REFERENCES category(name) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS posts_votes (
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    vote INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS comments_votes (
    comment_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    vote INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS posts_reports (
    id INTEGER NOT NULL PRIMARY KEY,
    post_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    moderator_id INTEGER NOT NULL,
    moderator_name TEXT NOT NULL,
    create_at DATETIME NOT NULL,
    FOREIGN KEY (moderator_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS notifications  (
    id INTEGER NOT NULL PRIMARY KEY,
    post_id INTEGER NOT NULL,
    comment_id INTEGER ,
    author_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    user_name TEXT NOT NULL,
    content TEXT NOT NULL,
    vote INTEGER NOT NULL,
    type INTEGER NOT NULL,
    create_at DATETIME NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
);
INSERT
    OR IGNORE INTO category (name)
VALUES ('Golang'),
    ('C++'),
    ('Java'),
    ('Python'),
    ('Kotlin'),
    ('Other');
INSERT
    OR IGNORE INTO users (name, email, password_hash, mode, rols)
VALUES (
        'admin',
        'admin@admin.com',
        '6d676664236735f4bff01e9b207dd4b0ed3433006025367b7c399f',
        0,
        10
    );