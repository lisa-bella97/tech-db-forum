package database

import (
	"github.com/lisa-bella97/tech-db-forum/app/models"
)

func Clear() {
	_, _ = Connection.Exec("TRUNCATE users, forums, threads, posts, votes, forum_users")
}

func GetStatus() *models.Status {
	status := models.Status{}
	_ = Connection.QueryRow(`SELECT (SELECT COUNT(*) FROM users) AS users, (SELECT COUNT(*) FROM forums) AS forums,
		(SELECT COUNT(*) FROM posts) AS posts, (SELECT COALESCE(SUM(threads), 0) FROM forums WHERE threads > 0)
		AS threads`).Scan(&status.User, &status.Forum, &status.Post, &status.Thread)
	return &status
}
