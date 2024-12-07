package usermodel

import (
	"database/sql"
	"go-web-native/config"
	"go-web-native/entities"
)

func GetByUsername(username string) (*entities.User, error) {
	var user entities.User
	row := config.DB.QueryRow("SELECT * FROM users WHERE username=?", username)
	if err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Kembalikan nil jika tidak ada user yang ditemukan
		}
		return nil, err // Kembalikan error jika terjadi kesalahan lain
	}
	return &user, nil // Kembalikan user yang ditemukan
}

func Create(user entities.User) error {
	_, err := config.DB.Exec(
		"INSERT INTO users (username, first_name, last_name, password) VALUES (?, ?, ?, ?)",
		user.Username, user.FirstName, user.LastName, user.Password,
	)
	if err != nil {
		return err
	}
	return nil
}