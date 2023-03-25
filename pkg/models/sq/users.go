package sq

import (
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
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
