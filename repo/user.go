package repo

import (
	"fmt"
	"swift_transit/domain"
	"swift_transit/user"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// UserRepo interface
type UserRepo interface {
	user.UserRepo
}

// userRepo struct
type userRepo struct {
	dbCon *sqlx.DB
}

// Constructor
func NewUserRepo(dbCon *sqlx.DB) UserRepo {
	return &userRepo{
		dbCon: dbCon,
	}
}

// Create new user with hashed password
func (r *userRepo) Create(user domain.User) (*domain.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO users (name, username, email, password, is_student, balance)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, username, email, is_student, balance
	`

	createdUser := domain.User{}
	err = r.dbCon.Get(
		&createdUser,
		query,
		user.Name,
		user.UserName,
		user.Email,
		string(hashedPassword),
		user.IsStudent,
		user.Balance,
	)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

// Find user by username and verify password (login)
func (r *userRepo) Find(userName, password string) (*domain.User, error) {
	user := domain.User{}
	query := `SELECT id, name, username, email, password, is_student, balance FROM users WHERE username=$1`

	err := r.dbCon.Get(&user, query, userName)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Remove password before returning
	user.Password = ""

	return &user, nil
}
