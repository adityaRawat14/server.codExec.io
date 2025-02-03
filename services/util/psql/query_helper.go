package psql

import (
	"context"
	"fmt"
	"server/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func GetUserByEmail(email string, db *pgxpool.Conn) (models.User, error) {
	var user models.User

	query := `SELECT id, first_name, last_name, email, image, created_at, updated_at FROM users WHERE email = $1`
	err := db.QueryRow(context.Background(), query, email).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Image,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {

		return models.User{}, err
	}

	return user, nil
}

func GetUserById(id uint, db *pgxpool.Conn) (models.User, error) {
	var user models.User
	query := `SELECT id, first_name, last_name, email, image, created_at, updated_at	FROM users WHERE id = $1`
	err := db.QueryRow(context.Background(), query, id).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Image,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, err
		}
		return models.User{}, err
	}

	return user, nil
}

func InsertUserIntoDb(user *models.NewUser, db *pgxpool.Conn) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hash error:", err)
	}

	query := `  INSERT INTO users (first_name, last_name, email, password,image ,created_at, updated_at)  VALUES ($1, $2, $3, $4,$5, NOW(), NOW())`
	_, err = db.Exec(context.Background(), query,
		user.FirstName,
		user.LastName,
		user.Email,
		hashedPassword,
		user.Image,
	)

	if err != nil {

		return err
	}

	return nil

}
