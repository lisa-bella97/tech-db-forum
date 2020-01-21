package database

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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

func CreatePosts(posts *models.Posts, thread *models.Thread) (*models.Posts, *models.ModelError) {
	dateTimeTemplate := "2006-01-02 15:04:05"
	created := time.Now().Format(dateTimeTemplate)
	query := strings.Builder{}
	query.WriteString(`INSERT INTO posts (parent, author, message, "isEdited", forum, thread, created, path) VALUES `)
	queryBody := `(%d, '%s', '%s', %t, '%s', %d, '%s', (SELECT path FROM posts WHERE id = %d) ||
		(SELECT last_value FROM posts_id_seq)),`
	postsNumber := len(*posts)

	for i, post := range *posts {
		e := CheckParent(post.Parent, thread.Id)
		if e != nil {
			return nil, e
		}
		temp := fmt.Sprintf(queryBody, post.Parent, post.Author, post.Message, post.IsEdited, thread.Forum, thread.Id,
			created, post.Parent)
		// удаление запятой в конце queryBody для последнего подзапроса
		if i == postsNumber-1 {
			temp = temp[:len(temp)-1]
		}
		query.WriteString(temp)
	}

	query.WriteString(" RETURNING author, created, forum, id, message, parent, thread")

	tx, txErr := Connection.Begin()
	if txErr != nil {
		return nil, &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   txErr.Error(),
		}
	}
	defer tx.Rollback()

	rows, err := tx.Query(query.String())
	defer rows.Close()
	if err != nil {
		return nil, &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}
	}

	insertPosts := models.Posts{}
	for rows.Next() {
		post := models.Post{}
		rows.Scan(
			&post.Author,
			&post.Created,
			&post.Forum,
			&post.Id,
			&post.Message,
			&post.Parent,
			&post.Thread,
		)
		insertPosts = append(insertPosts, post)
	}

	err = rows.Err()
	if err != nil {
		return nil, &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   err.Error(),
		}
	}

	tx.Commit()

	return &insertPosts, nil
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
