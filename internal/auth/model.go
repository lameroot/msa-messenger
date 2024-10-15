package auth

import (
	"database/sql"
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"

	"log"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidPassword = errors.New("invalid password: must be at least 6 characters long")
	ErrUserNotFound    = errors.New("user not found")
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Nickname     *string   `json:"nickname"`
	Info         *string   `json:"info"`
	AvatarUrl    *string   `json:"avatar_url"`
}

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) bool {
	return len(password) >= 6
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(db *sql.DB, email, password string) (*User, error) {
	if !ValidateEmail(email) {
		return nil, ErrInvalidEmail
	}

	if !ValidatePassword(password) {
		return nil, ErrInvalidPassword
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &User{
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	log.Printf("creating new user: %+v\n", user)
	err = db.QueryRow(
		"INSERT INTO users (email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
	//log.Printf("ERROR DB", err)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	user := &User{}
	err := db.QueryRow("SELECT id, email, password_hash, created_at, updated_at, nickname, avatar_url, info FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.Nickname, &user.AvatarUrl, &user.Info)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func GetUserByNickname(db *sql.DB, nickname string) (*User, error) {
	user := &User{}
	err := db.QueryRow("SELECT id, email, password_hash, created_at, updated_at, nickname, avatar_url, info FROM users WHERE nickname = $1", nickname).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.Nickname, &user.AvatarUrl, &user.Info)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func UpdateUser(db *sql.DB, id int64, nickname string, avatar_url string, info string) (*User, error) {
	user := &User{}
	err := db.QueryRow(
		`UPDATE users 
		SET nickname = $1, info = $2, avatar_url = $3 
		WHERE id = $4 
		RETURNING id, nickname, info, avatar_url`,
		nickname, info, avatar_url, id).
		Scan(&user.ID, &user.Nickname, &user.Info, &user.AvatarUrl)

	if err != nil {
		log.Printf("Error update: ", err)
		return nil, err
	}

	return user, nil
}
