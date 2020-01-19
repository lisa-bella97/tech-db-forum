package database

import (
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"net/http"
)

func GetForumBySlug(slug string) (models.Forum, *models.ModelError) {
	row := Connection.QueryRow(`SELECT * FROM forums WHERE LOWER(slug) = LOWER($1)`, slug)
	forum := models.Forum{}
	err := row.Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return models.Forum{}, &models.ModelError{
			ErrorCode: http.StatusNotFound,
			Message:   "Can't find forum with slug " + slug,
		}
	}
	return forum, nil
}

func CreateForum(forum models.Forum) *models.ModelError {
	_, err := Connection.Exec(`INSERT INTO forums (title, "user", slug) VALUES ($1, $2, $3)`,
		forum.Title, forum.User, forum.Slug)
	if err != nil {
		return &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Cannot create forum: " + err.Error(),
		}
	}
	return nil
}

func GetForumThreads(slug, limit, since string, desc bool) (models.Threads, *models.ModelError) {
	var result []models.Thread

	query := "SELECT * FROM threads WHERE forum = $1"
	if since != "" && desc {
		query += " AND created <= TIMESTAMPTZ '" + since + "'"
	} else if since != "" {
		query += " AND created >= TIMESTAMPTZ '" + since + "'"
	}
	query += " ORDER BY created"
	if desc {
		query += " DESC"
	}
	query += " LIMIT $2::TEXT::INTEGER"

	rows, err := Connection.Query(query, slug, limit)
	if err != nil {
		return []models.Thread{}, &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Cannot get forum threads: " + err.Error(),
		}
	}
	defer rows.Close()

	for rows.Next() {
		thread := models.Thread{}
		err = rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
			&thread.Slug, &thread.Created)
		if err != nil {
			return []models.Thread{}, &models.ModelError{
				ErrorCode: http.StatusInternalServerError,
				Message:   "Database query result parsing error: " + err.Error(),
			}
		}
		result = append(result, thread)
	}

	return result, nil
}

func UpdateForumPosts(slug string, posts int) *models.ModelError {
	_, err := Connection.Exec(`UPDATE forums SET posts = posts + $1 WHERE LOWER(slug) = LOWER($2)`, posts, slug)
	if err != nil {
		return &models.ModelError{
			ErrorCode: http.StatusNotFound,
			Message:   "Can't find forum with slug " + slug,
		}
	}
	return nil
}
