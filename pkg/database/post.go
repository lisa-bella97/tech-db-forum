package database

import (
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"net/http"
	"strconv"
)

func CreatePost(post *models.Post) *models.ModelError {
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

func GetPost(id int64) (*models.Post, *models.ModelError) {
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
	post, err := GetPost(id)
	postFull.Post = post
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
