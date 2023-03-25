package sq

import (
	"errors"
	"strings"

	"a4lab2.com/thoughtbin/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserModel struct {
	DB *gorm.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user := models.User{Name: name, Email: email, HashedPassword: hashedPass}
	err = m.DB.Create(&user).Error

	if err != nil {
		if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique") {
			return models.ErrDuplicateEmail
		}
	}
	return nil
}
func (m *UserModel) Authenticate(email, password string) (uint, error) {
	var user models.User
	// var id int
	// var hashedPass []byte
	row := m.DB.Select("ID", "HashedPassword").First(&user, "email = ? AND active = ?", email, true)

	if row.Error != nil {
		if errors.Is(row.Error, gorm.ErrRecordNotFound) {
			return 0, models.ErrInvalidCredentials
		}
	} else {
		return 0, row.Error
	}

	err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	//Otherwise, the password is correct. Return the user ID.
	return user.ID, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}
	row := m.DB.Select("ID", "name", "email", "created", "active").First(&u, "ID = ? ", id)
	if row.Error != nil {
		if errors.Is(row.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrNoRecord
		}
	} else {
		return nil, row.Error
	}
	return u, nil
}
