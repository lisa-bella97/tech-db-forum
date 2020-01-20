package database

import (
	"github.com/lisa-bella97/tech-db-forum/app/models"
	"github.com/pkg/errors"
	"net/http"
)

func GetUserByNickname(nickname string) (*models.User, *models.ModelError) {
	row := Connection.QueryRow(`SELECT * FROM users WHERE LOWER(nickname) = LOWER($1)`, nickname)
	user := models.User{}
	err := row.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		return nil, &models.ModelError{
			ErrorCode: http.StatusNotFound,
			Message:   "Can't find user with nickname " + nickname,
		}
	}
	return &user, nil
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

func GetUsersByNicknameOrEmail(nickname string, email string) (models.Users, *models.ModelError) {
	var result []models.User
	rows, err := Connection.Query(`SELECT * FROM users WHERE LOWER(nickname) = LOWER($1) OR LOWER(email) = LOWER($2)`,
		nickname, email)
	if err != nil {
		return []models.User{}, &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Cannot get users by nickname or email: " + err.Error(),
		}
	}
	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return []models.User{}, &models.ModelError{
				ErrorCode: http.StatusInternalServerError,
				Message:   "Database query result parsing error: " + err.Error(),
			}
		}
		result = append(result, user)
	}

	return result, nil
}

func CreateUser(user models.User) *models.ModelError {
	_, err := Connection.Exec(`INSERT INTO users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4)`,
		user.Nickname, user.Fullname, user.About, user.Email)
	if err != nil {
		return &models.ModelError{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Cannot create user: " + err.Error(),
		}
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
