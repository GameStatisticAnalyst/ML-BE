package models

import (
	"database/sql"
	"errors"
	"regexp"

	"github.com/GameStatisticAnalyst/ML-BE/db"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(username, email, password string) *User {
	return &User{
		ID:       uuid.NewString(),
		Username: username,
		Email:    email,
		Password: password,
	}
}

func (u *User) Save() error {
	db, err := db.NewDB()
	if err != nil {
		return err
	}

	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO user (id, username, email, password) VALUES (?, ?, ?, ?)", u.ID, u.Username, u.Email, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(email string) (*User, error) {
	db, err := db.NewDB()
	if err != nil {
		return nil, err
	}

	var user User
	query := db.QueryRow("SELECT * FROM user WHERE email = ?", email)
	err = query.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Article Not Found")
		}
		return nil, err
	}

	return &user, nil
}

func usernameValidation(username string) bool {
	return len(username) >= 3
}

func emailValidation(email string) bool {
	// Define the regex pattern for validating an email address
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex
	re := regexp.MustCompile(emailRegex)

	// Return true if the email matches the regex, false otherwise
	return re.MatchString(email)
}

func emailExist(email string) (bool, error) {
	db, err := db.NewDB()
	if err != nil {
		return false, err
	}
	defer db.Close()

	// Prepare the SQL query
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM user WHERE email = ?)`
	err = db.QueryRow(query, email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists, nil
}

func (u *User) ValidateRegisterUser() (bool, map[string]error) {
	errs := make(map[string]error)
	isValid := true

	if !usernameValidation(u.Username) {
		errs["username"] = errors.New("Username should contain atleast 3 characters")
		isValid = false
	}

	if !emailValidation(u.Email) {
		errs["email"] = errors.New("Email is not valid")
		isValid = false
	}

	email_exist, err := emailExist(u.Email)
	if err != nil {
		return false, nil
	}

	if email_exist {
		errs["email"] = errors.New("Email already been used")
		isValid = false
	}

	if len(u.Password) < 6 {
		errs["password"] = errors.New("Password should atleast contain 6 characters")
		isValid = false
	}

	return isValid, errs
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
