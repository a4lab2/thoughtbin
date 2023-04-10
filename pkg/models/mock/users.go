package mock

import (
	"time"

	"a4lab2.com/thoughtbin/pkg/models"
)

var mockUser = &models.User{
	Name:    "Alice",
	Email:   "alice@example.com",
	Created: time.Now(),
	Active:  true,
}

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}
func (m *UserModel) Authenticate(email, password string) (uint, error) {
	switch email {
	case "alice@example.com":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}
func (m *UserModel) Get(id uint) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}
