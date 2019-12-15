package database

import (
	"github.com/jackc/pgx"
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"github.com/pkg/errors"
)

var Connection *pgx.ConnPool

func Init() {
	Connection, _ = pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host: "localhost",
				Port: 5432,
				Database: "forum",
				User:     "forum",
				Password: "forum",
			},
			MaxConnections: 50,
		})
}

func GetUserByNickname(nickname string) (models.User, error) {
	res, err := Connection.Query(`SELECT * FROM users WHERE nickname = $1`, nickname)
	if err != nil {
		return models.User{}, errors.Wrap(err, "cannot get user by nickname")
	}
	defer res.Close()

	u := models.User{}

	if res.Next() {
		err = res.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
		if err != nil {
			return models.User{}, errors.Wrap(err, "db query result parsing error")
		}
		return u, nil
	}
	return models.User{}, errors.New("cannot get user by nickname")
}
