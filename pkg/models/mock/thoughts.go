package mock

import (
	"time"

	"a4lab2.com/thoughtbin/pkg/models"
	"gorm.io/gorm"
)

var mockThought = &models.Thought{
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type ThoughtModel struct {
	DB *gorm.DB
}

func (m *ThoughtModel) Insert(title, content, expires string) (uint, error) {
	return 2, nil
}

func (m *ThoughtModel) Get(id uint) (*models.Thought, error) {
	switch id {
	case 1:
		return mockThought, nil
	default:
		return nil, models.ErrNoRecord
	}
}
func (m *ThoughtModel) Latest() ([]*models.Thought, error) {
	return []*models.Thought{mockThought}, nil
}
