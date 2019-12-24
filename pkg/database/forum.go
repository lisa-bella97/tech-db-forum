package database

import (
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"github.com/pkg/errors"
)

func GetForumBySlug(slug string) (models.Forum, error) {
	rows, err := Connection.Query(`SELECT * FROM forums WHERE slug = $1`, slug)
	if err != nil {
		return models.Forum{}, errors.Wrap(err, "cannot get forum by slug")
	}
	defer rows.Close()

	if rows.Next() {
		forum := models.Forum{}
		err := rows.Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
		if err != nil {
			return models.Forum{}, errors.Wrap(err, "db query result parsing error")
		}
		return forum, nil
	}

	return models.Forum{}, errors.New("cannot find forum by slug")
}

func CreateForum(forum models.Forum) error {
	_, err := Connection.Exec(`INSERT INTO forums (title, "user", slug, posts, threads) VALUES ($1, $2, $3, $4, $5)`,
		forum.Title, forum.User, forum.Slug, forum.Posts, forum.Threads)
	if err != nil {
		return errors.Wrap(err, "cannot create forum")
	}

	return nil
}
