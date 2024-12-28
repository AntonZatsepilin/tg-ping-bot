package repository

import (
	"goPingRobot/auth/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
    query := "INSERT INTO users (username, password_hash) VALUES ($1, $2)"
    _, err := r.db.Exec(query, user.Username, user.PasswordHash)
    return err
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
    var user models.User
    query := "SELECT * FROM users WHERE username = $1"
    err := r.db.Get(&user, query, username)
    return &user, err
}