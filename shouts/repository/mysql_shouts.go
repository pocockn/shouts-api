package repository

import (
	"github.com/jinzhu/gorm"
	shoutsModel "github.com/pocockn/models/api/shouts"
	"github.com/pocockn/shouts-api/shouts"
)

type shoutsRepository struct {
	Conn *gorm.DB
}

// NewShoutsRepository creates a new shoutRepository struct.
func NewShoutsRepository(conn *gorm.DB) shouts.Repository {
	return &shoutsRepository{conn}
}

// Fetch fetches a shout via it's ID from the DB.
func (s *shoutsRepository) Fetch(id uint) (shoutsModel.Shout, error) {
	var shout shoutsModel.Shout
	err := s.Conn.Where("id = ?", id).First(&shout).Error

	return shout, err
}

// FetchAll fetches all the shouts from the DB.
func (s *shoutsRepository) FetchAll() ([]shoutsModel.Shout, error) {
	var shouts []shoutsModel.Shout
	err := s.Conn.Find(&shouts).Error

	return shouts, err
}

// Store stores a shout in the DB.
func (s *shoutsRepository) Store(shout *shoutsModel.Shout) error {
	err := s.Conn.Create(shout).Error
	return err
}
