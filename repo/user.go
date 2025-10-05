package repo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	Id        int64   `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	UserName  string  `json:"username" db:"username"`
	Email     string  `json:"email" db:"email"`
	Password  string  `json:"password" db:"password"`
	IsStudent bool    `json:"is_student" db:"is_student"`
	Balance   float32 `json:"balance" db:"balance"`
}

// UserRepo interface
type UserRepo interface {
	Find(userName, password string) (*User, error) // login
	Create(user User) (*User, error)               // create new user
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
func (r *userRepo) Create(user User) (*User, error) {
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

	createdUser := User{}
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
func (r *userRepo) Find(userName, password string) (*User, error) {
	user := User{}
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
