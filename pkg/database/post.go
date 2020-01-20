package database

import (
	"github.com/jackc/pgx"
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"net/http"
	"strconv"
)

const postID = `
	SELECT id
	FROM posts
	WHERE id = $1 AND thread IN (SELECT id FROM threads WHERE thread <> $2)
`

func parentExitsInOtherThread(parent int64, threadID int32) bool {
	var id int64
	err := Connection.QueryRow(`SELECT id FROM posts WHERE id = $1 AND
		thread IN (SELECT id FROM threads WHERE thread <> $2)`, parent, threadID).Scan(&id)

	return err == nil || err != pgx.ErrNoRows
}

func parentNotExists(parent int64) bool {
	if parent == 0 {
		return false
	}

	var id int64
	err := Connection.QueryRow("SELECT id FROM posts WHERE id = $1", parent).Scan(&id)

	return err != nil
}

func CheckParent(parent int64, threadId int32) *models.ModelError {
	if parentExitsInOtherThread(parent, threadId) || parentNotExists(parent) {
		return &models.ModelError{
			ErrorCode: http.StatusConflict,
			Message:   "Post parent not found",
		}
	}
	return nil
}

func CreatePost(post *models.Post) *models.ModelError {
	e := CheckParent(post.Parent, post.Thread)
	if e != nil {
		return e
	}

	err := Connection.QueryRow(`INSERT INTO posts (parent, author, message, "isEdited", forum, thread, created, path)
		VALUES ($1, $2, $3, $4, $5, $6, $7, (SELECT path FROM posts WHERE id = $1) ||
		(SELECT last_value FROM posts_id_seq)) RETURNING id, created`, post.Parent, post.Author, post.Message,
		post.IsEdited, post.Forum, post.Thread, post.Created).Scan(&post.Id, &post.Created)
	if err != nil {
		return &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Cannot create post: " + err.Error(),
		}
	}
	return nil
}

func UpdateForumUsers(posts []models.Post) {
	for _, p := range posts {
		_, _ = Connection.Exec(`INSERT INTO forum_users VALUES ($1, $2) ON CONFLICT DO NOTHING`, p.Author, p.Forum)
	}
}

func GetPostById(id int64) (*models.Post, *models.ModelError) {
	post := models.Post{}
	err := Connection.QueryRow(`SELECT id, parent, author, message, "isEdited", forum, thread, created FROM posts 
			WHERE id = $1`, id).Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum,
		&post.Thread, &post.Created)
	if err != nil {
		return nil, &models.ModelError{
			ErrorCode: http.StatusNotFound,
			Message:   "Cannot get post with ID " + strconv.Itoa(int(id)),
		}
	}
	return &post, nil
}

func GetPostFull(id int64, related []string) (*models.PostFull, *models.ModelError) {
	postFull := models.PostFull{}
	var err *models.ModelError
	postFull.Post, err = GetPostById(id)
	if err != nil {
		return nil, err
	}

	for _, obj := range related {
		switch obj {
		case "user":
			postFull.Author, err = GetUserByNickname(postFull.Post.Author)
		case "forum":
			postFull.Forum, err = GetForumBySlug(postFull.Post.Forum)
		case "thread":
			postFull.Thread, err = GetThreadById(postFull.Post.Thread)
		}

		if err != nil {
			return nil, err
		}
	}

	return &postFull, nil
}

func UpdatePost(postUpdate *models.PostUpdate, id int64) (*models.Post, *models.ModelError) {
	post, err := GetPostById(id)
	if err != nil {
		return nil, err
	}

	if len(postUpdate.Message) == 0 {
		return post, nil
	}

	e := Connection.QueryRow(`UPDATE posts SET message = COALESCE($1, message),
		"isEdited" = ($1 IS NOT NULL AND $1 <> message) WHERE id = $2
		RETURNING id, parent, author, message, "isEdited", forum, thread, created`, &postUpdate.Message, id).Scan(&post.Id,
		&post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	if e != nil {
		return nil, &models.ModelError{
			ErrorCode: http.StatusNotFound,
			Message:   "Cannot get post with ID " + strconv.Itoa(int(id)),
		}
	}

	return post, nil
}
