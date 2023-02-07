package sq

import (
	"time"

	"a4lab2.com/thoughtbin/pkg/models"
	"gorm.io/gorm"
)

type ThoughtModel struct {
	DB *gorm.DB
}

// title, content, expires string
func (m *ThoughtModel) Insert(title, content, expires string) (uint, error) {

	thought := models.Thought{Title: title, Content: content, Expires: time.Now()}
	_ = m.DB.Create(&thought)
	id := thought.ID
	return id, nil
}

func (m *ThoughtModel) BatchInsert(thoughts []models.Thought) ([]uint, error) {
	ids := make([]uint, 5)
	_ = m.DB.Create(&thoughts)
	for _, thought := range thoughts {
		ids = append(ids, thought.ID)
	}

	return ids, nil
}

func (m *ThoughtModel) Get(id int) (*models.Thought, error) {
	var thought models.Thought
	err := m.DB.First(&thought, id).Error // Check if returns RecordNotFound error

	return &thought, err
}

func (m *ThoughtModel) Latest() ([]*models.Thought, error) {
	var thoughts []*models.Thought
	err := m.DB.Order("ID desc").Limit(3).Find(&thoughts).Error
	return thoughts, err
}