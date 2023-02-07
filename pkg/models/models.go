package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var ErrNoRecord = errors.New("models:no matching record")

type Thought struct {
	gorm.Model
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
