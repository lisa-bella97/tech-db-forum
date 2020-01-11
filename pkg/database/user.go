package database

import (
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"github.com/pkg/errors"
)

func GetUserByNickname(nickname string) (models.User, error) {
	rows, err := Connection.Query(`SELECT * FROM users WHERE nickname = $1`, nickname)
	if err != nil {
		return models.User{}, errors.Wrap(err, "cannot get user by nickname")
	}
	defer rows.Close()

	if rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return models.User{}, errors.Wrap(err, "db query result parsing error")
		}
		return user, nil
	}

	return models.User{}, errors.New("user not found by nickname")
}

func GetUsersByNicknameOrEmail(nickname string, email string) (models.Users, error) {
	var result []models.User
	rows, err := Connection.Query(`SELECT * FROM users WHERE LOWER(nickname) = LOWER($1) OR LOWER(email) = LOWER($2)`,
		nickname, email)
	if err != nil {
		return []models.User{}, errors.Wrap(err, "cannot get users by nickname or email")
	}
	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return []models.User{}, errors.Wrap(err, "db query result parsing error")
		}
		result = append(result, user)
	}

	return result, nil
}

func CreateUser(user models.User) error {
	_, err := Connection.Exec(`INSERT INTO users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4)`,
		user.Nickname, user.Fullname, user.About, user.Email)
	if err != nil {
		return errors.Wrap(err, "cannot create user")
	}

	return nil
}
