DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "forums" CASCADE;
DROP TABLE IF EXISTS "threads" CASCADE;
DROP TABLE IF EXISTS "posts" CASCADE;
DROP TABLE IF EXISTS "votes" CASCADE;
DROP TABLE IF EXISTS "forum_users" CASCADE;

CREATE TABLE users
(
    "nickname" TEXT UNIQUE PRIMARY KEY,
    "fullname" TEXT        NOT NULL,
    "about"    TEXT,
    "email"    TEXT UNIQUE NOT NULL
);

CREATE TABLE forums
(
    "title"   TEXT        NOT NULL,
    "user"    TEXT        NOT NULL REFERENCES users ("nickname"),
    "slug"    TEXT UNIQUE NOT NULL,
    "posts"   BIGINT  DEFAULT 0,
    "threads" INTEGER DEFAULT 0
);

CREATE TABLE threads
(
    "id"      SERIAL UNIQUE PRIMARY KEY,
    "title"   TEXT NOT NULL,
    "author"  TEXT NOT NULL REFERENCES users ("nickname"),
    "forum"   TEXT NOT NULL REFERENCES forums ("slug"),
    "message" TEXT NOT NULL,
    "votes"   INTEGER                  DEFAULT 0,
    "slug"    TEXT,
    "created" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE posts
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

CREATE TABLE votes
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

DROP INDEX IF EXISTS idx_users_nickname;
DROP INDEX IF EXISTS idx_users_nickname_email;
DROP INDEX IF EXISTS idx_forums_slug;
DROP INDEX IF EXISTS idx_threads_id;
DROP INDEX IF EXISTS idx_threads_slug;
DROP INDEX IF EXISTS idx_threads_created_forum;
DROP INDEX IF EXISTS idx_posts_id;
DROP INDEX IF EXISTS idx_posts_thread_id;
DROP INDEX IF EXISTS idx_posts_thread_id0;
DROP INDEX IF EXISTS idx_posts_thread_path1_id;
DROP INDEX IF EXISTS idx_posts_thread_path_parent;
DROP INDEX IF EXISTS idx_posts_thread;
DROP INDEX IF EXISTS idx_posts_path_AA;
DROP INDEX IF EXISTS idx_posts_path_AD;
DROP INDEX IF EXISTS idx_posts_path_DA;
DROP INDEX IF EXISTS idx_posts_path_DD;
DROP INDEX IF EXISTS idx_posts_path_desc;
DROP INDEX IF EXISTS idx_posts_paths;
DROP INDEX IF EXISTS idx_posts_thread_path;
DROP INDEX IF EXISTS idx_posts_thread_id_created;
DROP INDEX IF EXISTS idx_votes_thread_nickname;

DROP INDEX IF EXISTS idx_fu_user;
DROP INDEX IF EXISTS idx_fu_forum;

CREATE INDEX IF NOT EXISTS idx_fu_user ON forum_users (forum, forum_user);
CREATE INDEX IF NOT EXISTS idx_fu_forum ON forum_users (forum);

CREATE INDEX IF NOT EXISTS idx_users_nickname ON users (nickname);

CREATE INDEX IF NOT EXISTS idx_forums_slug ON forums (slug);

CREATE INDEX IF NOT EXISTS idx_threads_id ON threads (id);
CREATE INDEX IF NOT EXISTS idx_threads_slug ON threads (slug);
CREATE INDEX IF NOT EXISTS idx_threads_forum ON threads (forum);

CREATE INDEX IF NOT EXISTS idx_posts_forum ON posts (forum);
CREATE INDEX IF NOT EXISTS idx_posts_id ON posts (id);
CREATE INDEX IF NOT EXISTS idx_posts_thread_id ON posts (thread, id);
CREATE INDEX IF NOT EXISTS idx_posts_thread_id0 ON posts (thread, id) WHERE parent = 0;
CREATE INDEX IF NOT EXISTS idx_posts_thread_id_created ON posts (id, created, thread);

CREATE UNIQUE INDEX IF NOT EXISTS idx_votes_thread_nickname ON votes (thread, nickname);


DROP FUNCTION IF EXISTS insert_vote();
CREATE OR REPLACE FUNCTION insert_vote() RETURNS TRIGGER AS
$insert_vote$
BEGIN
    UPDATE threads
    SET votes = votes + NEW.voice
    WHERE id = NEW.thread;
    RETURN NEW;
END;
$insert_vote$
    LANGUAGE plpgsql;
DROP TRIGGER IF EXISTS insert_vote ON votes;
CREATE TRIGGER insert_vote
    BEFORE INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE insert_vote();


DROP FUNCTION IF EXISTS update_vote();
CREATE OR REPLACE FUNCTION update_vote() RETURNS TRIGGER AS
$update_vote$
BEGIN
    UPDATE threads
    SET votes = votes - OLD.voice + NEW.voice
    WHERE id = NEW.thread;
    RETURN NEW;
END;
$update_vote$
    LANGUAGE plpgsql;
DROP TRIGGER IF EXISTS update_vote ON votes;
CREATE TRIGGER update_vote
    BEFORE UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE update_vote();


DROP FUNCTION IF EXISTS add_forum_user();
CREATE OR REPLACE FUNCTION add_forum_user() RETURNS TRIGGER AS
$add_forum_user$
BEGIN
    INSERT INTO forum_users VALUES (NEW.author, NEW.forum) ON CONFLICT DO NOTHING;
    RETURN NULL;
END;
$add_forum_user$
    LANGUAGE plpgsql;
CREATE TRIGGER add_forum_user
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE add_forum_user();
