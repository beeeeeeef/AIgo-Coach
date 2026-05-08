package repository

import (
	"aigo-coach/backend/model"
    "database/sql"    
)

type UserRepositorySQL struct {
	DB *sql.DB
}
func NewUserRepositorySQL(db *sql.DB) *UserRepositorySQL {
	return &UserRepositorySQL{DB: db}
}
func (r *UserRepositorySQL) GetByEmail(email string) ( *model.User,error){
	query := `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`

	var user model.User

	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil


}
func (r *UserRepositorySQL) GetByID(id int64) ( *model.User,error){
	query := `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE id = ?
		LIMIT 1
	`

	var user model.User

	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil

}
func (r *UserRepositorySQL) Create(user *model.User) error{
		query := `
		INSERT INTO users (email, password_hash)
		VALUES (?, ?)
	`

	result, err := r.DB.Exec(query, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil

}
