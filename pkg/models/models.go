package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrNoRecord           = errors.New("models:no matching record")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type Thought struct {
	gorm.Model
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	gorm.Model
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}
