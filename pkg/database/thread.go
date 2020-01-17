package database

import (
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"net/http"
	"strconv"
)

func GetThreadBySlug(slug string) (models.Thread, *models.ModelError) {
	rows, err := Connection.Query(`SELECT * FROM threads WHERE LOWER(slug) = LOWER($1)`, slug)
	if err != nil {
		return models.Thread{}, &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Cannot get thread by slug: " + err.Error(),
		}
	}
	defer rows.Close()

	if rows.Next() {
		thread := models.Thread{}
		err := rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message,
			&thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return models.Thread{}, &models.ModelError{
				ErrorCode: http.StatusInternalServerError,
				Message:   "Database query result parsing error: " + err.Error(),
			}
		}
		return thread, nil
	}

	return models.Thread{}, &models.ModelError{
		ErrorCode: http.StatusNotFound,
		Message:   "Can't find thread with slug "+slug,
	}
}

func GetThreadById(id int) (models.Thread, *models.ModelError) {
	rows, err := Connection.Query(`SELECT * FROM threads WHERE id = $1`, id)
	if err != nil {
		return models.Thread{}, &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Cannot get thread by ID: " + err.Error(),
		}
	}
	defer rows.Close()

	if rows.Next() {
		thread := models.Thread{}
		err := rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message,
			&thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return models.Thread{}, &models.ModelError{
				ErrorCode: http.StatusInternalServerError,
				Message:   "Database query result parsing error: " + err.Error(),
			}
		}
		return thread, nil
	}

	return models.Thread{}, &models.ModelError{
		ErrorCode: http.StatusNotFound,
		Message:   "Can't find thread with ID " + strconv.Itoa(id),
	}
}

func CreateThread(thread *models.Thread) *models.ModelError {
	err := Connection.QueryRow(`INSERT INTO threads (title, author, forum, message, slug, created)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, thread.Title, thread.Author, thread.Forum, thread.Message,
		thread.Slug, thread.Created).Scan(&thread.Id)
	if err != nil {
		return &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Cannot create thread: " + err.Error(),
		}
	}
	return nil
}

func Vote(vote models.Vote, threadId int32) *models.ModelError {
	_, err := Connection.Exec(`INSERT INTO votes VALUES ($1, $2, $3)`, vote.Nickname, vote.Voice, threadId)
	if err != nil {
		_, err = Connection.Exec(`UPDATE votes SET voice = $1 WHERE thread = $2 AND nickname = $3`,
			vote.Voice, threadId, vote.Nickname)
		if err != nil {
			return &models.ModelError{
				ErrorCode: http.StatusNotFound,
				Message:   "Can't find user with nickname " + vote.Nickname,
			}
		}
	}

	_, err = Connection.Exec(`UPDATE threads SET votes = votes+$1 WHERE id = $2`, vote.Voice, threadId)
	if err != nil {
		return &models.ModelError{
			ErrorCode: http.StatusNotFound,
			Message:   "Can't find thread with ID " + strconv.Itoa(int(threadId)),
		}
	}

	return nil
}
