DROP TABLE IF EXISTS "forums" CASCADE;
DROP TABLE IF EXISTS "posts" CASCADE;
DROP TABLE IF EXISTS "threads" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "votes" CASCADE;
DROP TABLE IF EXISTS "forum_users" CASCADE;

CREATE TABLE IF NOT EXISTS users
(
    "nickname" TEXT UNIQUE PRIMARY KEY,
    "fullname" TEXT        NOT NULL,
    "about"    TEXT,
    "email"    TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS forums
(
    "title"   TEXT        NOT NULL,
    "user"    TEXT        NOT NULL REFERENCES users ("nickname"),
    "slug"    TEXT UNIQUE NOT NULL,
    "posts"   BIGINT  DEFAULT 0,
    "threads" INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS threads
(
    "id"      SERIAL UNIQUE PRIMARY KEY,
    "title"   TEXT NOT NULL,
    "author"  TEXT NOT NULL REFERENCES users ("nickname"),
    "forum"   TEXT NOT NULL REFERENCES forums ("slug"),
    "message" TEXT NOT NULL,
    "votes"   INTEGER                  DEFAULT 0,
    "slug"    TEXT                     DEFAULT NULL UNIQUE,
    "created" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts
(
    "id"       BIGSERIAL UNIQUE PRIMARY KEY,
    "parent"   INTEGER                  DEFAULT 0,
    "author"   TEXT    NOT NULL REFERENCES users ("nickname"),
    "message"  TEXT    NOT NULL,
    "isEdited" BOOLEAN                  DEFAULT FALSE,
    "forum"    TEXT    NOT NULL REFERENCES forums ("slug"),
    "thread"   INTEGER NOT NULL REFERENCES threads ("id"),
    "created"  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "path"     BIGINT[]
);

CREATE TABLE IF NOT EXISTS votes
(
    "nickname" TEXT    NOT NULL,
    "voice"    INTEGER NOT NULL,
    "thread"   INT     NOT NULL REFERENCES threads ("id")
);


CREATE TABLE forum_users
(
    "forum_user" TEXT COLLATE ucs_basic NOT NULL,
    "forum"      TEXT                   NOT NULL
);
