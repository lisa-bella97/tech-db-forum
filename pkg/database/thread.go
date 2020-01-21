package database

import (
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"net/http"
	"strconv"
)

func GetThreadBySlug(slug string) (*models.Thread, *models.ModelError) {
	rows, err := Connection.Query(`SELECT * FROM threads WHERE LOWER(slug) = LOWER($1)`, slug)
	if err != nil {
		return nil, &models.ModelError{
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
			return nil, &models.ModelError{
				ErrorCode: http.StatusInternalServerError,
				Message:   "Database query result parsing error: " + err.Error(),
			}
		}
		return &thread, nil
	}

	return nil, &models.ModelError{
		ErrorCode: http.StatusNotFound,
		Message:   "Can't find thread with slug " + slug,
	}
}

func GetThreadById(id int32) (*models.Thread, *models.ModelError) {
	rows, err := Connection.Query(`SELECT * FROM threads WHERE id = $1`, id)
	if err != nil {
		return nil, &models.ModelError{
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
			return nil, &models.ModelError{
				ErrorCode: http.StatusInternalServerError,
				Message:   "Database query result parsing error: " + err.Error(),
			}
		}
		return &thread, nil
	}

	return nil, &models.ModelError{
		ErrorCode: http.StatusNotFound,
		Message:   "Can't find thread with ID " + strconv.Itoa(int(id)),
	}
}

func isNumber(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func GetThread(param string) (*models.Thread, *models.ModelError) {
	var err *models.ModelError
	var thread *models.Thread

	if isNumber(param) {
		id, _ := strconv.Atoi(param)
		thread, err = GetThreadById(int32(id))
	} else {
		thread, err = GetThreadBySlug(param)
	}

	if err != nil {
		return nil, err
	}

	return thread, nil
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

// Возвращает новую оценку ветви обсуждения
func Vote(vote models.Vote, threadId int32) (int32, *models.ModelError) {
	_, err := Connection.Exec(`INSERT INTO votes VALUES ($1, $2, $3)`, vote.Nickname, vote.Voice, threadId)
	if err != nil {
		_, err = Connection.Exec(`UPDATE votes SET voice = $1 WHERE thread = $2 AND nickname = $3`,
			vote.Voice, threadId, vote.Nickname)
		if err != nil {
			return 0, &models.ModelError{
				ErrorCode: http.StatusNotFound,
				Message:   "Can't find user with nickname " + vote.Nickname,
			}
		}
	}

	var newVotes int32

	err = Connection.QueryRow("SELECT votes FROM threads WHERE id = $1", threadId).Scan(&newVotes)
	if err != nil {
		return 0, &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Can't find votes from thread with ID " + strconv.Itoa(int(threadId)) + ": " + err.Error(),
		}
	}

	return newVotes, nil
}

func FlatSort(threadId int32, limit, since int, desc bool) string {
	query := `SELECT id, parent, author, message, "isEdited", forum, thread, created FROM posts WHERE thread = ` +
		strconv.Itoa(int(threadId))

	if since != 0 {
		if desc {
			query += " AND id < " + strconv.Itoa(since)
		} else {
			query += " AND id > " + strconv.Itoa(since)
		}
	}

	if desc {
		query += " ORDER BY id DESC"
	} else {
		query += " ORDER BY id"
	}

	query += " LIMIT " + strconv.Itoa(limit)

	return query
}

func TreeSort(threadId int32, limit, since int, desc bool) string {
	query := `SELECT id, parent, author, message, "isEdited", forum, thread, created FROM posts WHERE thread = ` +
		strconv.Itoa(int(threadId))

	if since != 0 {
		if desc {
			query += " AND path < (SELECT path FROM posts WHERE id = " + strconv.Itoa(since) + ")"
		} else {
			query += " AND path > (SELECT path FROM posts WHERE id = " + strconv.Itoa(since) + ")"
		}
	}

	if desc {
		query += " ORDER BY path DESC"
	} else {
		query += " ORDER BY path"
	}

	query += " LIMIT " + strconv.Itoa(limit)

	return query
}

func ParentTreeSort(threadId int32, limit, since int, desc bool) string {
	query := `SELECT id, parent, author, message, "isEdited", forum, thread, created FROM posts p WHERE
		p.thread = ` + strconv.Itoa(int(threadId)) + ` AND path[1] IN (SELECT p2.path[1] FROM posts p2 WHERE
		p2.thread = ` + strconv.Itoa(int(threadId)) + " AND p2.parent = 0"

	if since != 0 {
		if desc {
			query += " AND p2.path[1] < (SELECT p3.path[1] FROM posts p3 WHERE p3.id = " + strconv.Itoa(since) + ")"
		} else {
			query += " AND p2.path[1] > (SELECT p3.path[1] FROM posts p3 WHERE p3.id = " + strconv.Itoa(since) + ")"
		}
	}

	if desc {
		query += " ORDER BY p2.path DESC"
	} else {
		query += " ORDER BY p2.path"
	}

	query += " LIMIT " + strconv.Itoa(limit) + ")"

	if desc {
		query += " ORDER BY p.path[1] DESC, p.path[2:]"
	} else {
		query += " ORDER BY p.path"
	}

	return query
}

func GetThreadPosts(threadId int32, limit, since, sort string, desc bool) (models.Posts, *models.ModelError) {
	var result []models.Post
	query := ""
	limitInt, _ := strconv.Atoi(limit)
	sinceInt, _ := strconv.Atoi(since)

	switch sort {
	case "flat":
		query = FlatSort(threadId, limitInt, sinceInt, desc)
	case "tree":
		query = TreeSort(threadId, limitInt, sinceInt, desc)
	case "parent_tree":
		query = ParentTreeSort(threadId, limitInt, sinceInt, desc)
	}

	rows, err := Connection.Query(query)
	if err != nil {
		return []models.Post{}, &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Cannot get thread posts: " + err.Error(),
		}
	}
	defer rows.Close()

	for rows.Next() {
		post := models.Post{}
		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread,
			&post.Created)
		if err != nil {
			return []models.Post{}, &models.ModelError{
				ErrorCode: http.StatusInternalServerError,
				Message:   "Database query result parsing error: " + err.Error(),
			}
		}
		result = append(result, post)
	}

	return result, nil
}

func UpdateThread(thread *models.Thread) *models.ModelError {
	row := Connection.QueryRow(`UPDATE threads SET title = COALESCE(NULLIF($1, ''), title),
		message  = COALESCE(NULLIF($2, ''), message) WHERE LOWER(slug) = LOWER($3)
		RETURNING id, title, author, forum, message, votes, slug, created`, &thread.Title, &thread.Message,
		&thread.Slug)
	err := row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
		&thread.Slug, &thread.Created)
	if err != nil {
		return &models.ModelError{
			ErrorCode: http.StatusNotFound,
			Message:   "Can't find thread with slug " + thread.Slug,
		}
	}
	return nil
}
