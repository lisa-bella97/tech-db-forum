DROP TABLE IF EXISTS "forums" CASCADE;
DROP TABLE IF EXISTS "posts" CASCADE;
DROP TABLE IF EXISTS "threads" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "votes" CASCADE;
DROP TABLE IF EXISTS "forum_users" CASCADE;

CREATE TABLE IF NOT EXISTS users
(
    "nickname" TEXT UNIQUE PRIMARY KEY,
    "email"    TEXT UNIQUE NOT NULL,
    "fullname" TEXT        NOT NULL,
    "about"    TEXT
);

CREATE TABLE IF NOT EXISTS forums
(
    "posts"   BIGINT  DEFAULT 0,
    "slug"    TEXT UNIQUE NOT NULL,
    "threads" INTEGER DEFAULT 0,
    "title"   TEXT        NOT NULL,
    "user"    TEXT        NOT NULL REFERENCES users ("nickname")
);

CREATE TABLE IF NOT EXISTS threads
(
    "id"      SERIAL UNIQUE PRIMARY KEY,
    "author"  TEXT NOT NULL REFERENCES users ("nickname"),
    "created" TIMESTAMPTZ(3) DEFAULT now(),
    "forum"   TEXT NOT NULL REFERENCES forums ("slug"),
    "message" TEXT NOT NULL,
    "slug"    TEXT,
    "title"   TEXT NOT NULL,
    "votes"   INTEGER        DEFAULT 0
);

CREATE TABLE IF NOT EXISTS posts
(
    "id"       BIGSERIAL UNIQUE PRIMARY KEY,
    "author"   TEXT    NOT NULL REFERENCES users ("nickname"),
    "created"  TIMESTAMPTZ(3) DEFAULT now(),
    "forum"    TEXT    NOT NULL REFERENCES forums ("slug"),
    "isEdited" BOOLEAN        DEFAULT FALSE,
    "message"  TEXT    NOT NULL,
    "parent"   INTEGER        DEFAULT 0,
    "thread"   INTEGER NOT NULL REFERENCES threads ("id"),
    "path"     BIGINT[]
);

CREATE TABLE IF NOT EXISTS votes
(
    "thread"   INT     NOT NULL REFERENCES threads ("id"),
    "voice"    INTEGER NOT NULL,
    "nickname" TEXT    NOT NULL
);


CREATE TABLE forum_users
(
    "forum_user" TEXT COLLATE ucs_basic NOT NULL,
    "forum"      TEXT                   NOT NULL
);