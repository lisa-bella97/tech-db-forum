package database

import (
	"database/sql"
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"github.com/pkg/errors"
	"net/http"
)

func GetUserByNickname(nickname string) (models.User, *models.ModelError) {
	row := Connection.QueryRow(`SELECT * FROM users WHERE LOWER(nickname) = LOWER($1)`, nickname)
	user := models.User{}
	err := row.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err == sql.ErrNoRows {
		return models.User{}, &models.ModelError{
			ErrorCode: http.StatusNotFound,
			Message:   "Can't find user with nickname " + nickname,
		}
	} else if err != nil {
		return models.User{}, &models.ModelError{
			ErrorCode: http.StatusNotFound,
			Message:   "db query result parsing error: " + err.Error(),
		}
	}
	return user, nil
}

func GetUserByEmail(email string) (models.User, error) {
	row := Connection.QueryRow(`SELECT * FROM users WHERE LOWER(email) = LOWER($1)`, email)
	user := models.User{}
	err := row.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		return models.User{}, errors.New("Can't find user with email " + email)
	}
	return user, nil
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

func UpdateUser(user *models.User) *models.ModelError {
	row := Connection.QueryRow(`UPDATE users SET fullname = COALESCE(NULLIF($2, ''), fullname),
		about = COALESCE(NULLIF($3, ''), about), email = COALESCE(NULLIF($4, ''), email)
		WHERE nickname = $1 RETURNING nickname, fullname, about, email`,
		&user.Nickname, &user.Fullname, &user.About, &user.Email)
	err := row.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		return &models.ModelError{
			ErrorCode: http.StatusNotFound,
			Message:   "Can't find user with nickname " + user.Nickname,
		}
	}
	return nil
}
