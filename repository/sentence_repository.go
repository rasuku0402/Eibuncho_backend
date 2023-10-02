package repository

import (
	"fmt"
	"testToDoRestAPI/model"

	"gorm.io/gorm"
)

type ISentenceRepository interface {
	GetAllSentences(sentences *[]model.Sentences, userId uint) error
	GetSentenceById(sentence *model.Sentences, userId uint, sentenceId uint) error
	CreateSentence(sentence *model.Sentences) error
	DeleteSentence(userId uint, sentenceId uint) error
}

type sentenceRepository struct {
	db *gorm.DB
}

func NewSentenceRepository(db *gorm.DB) ISentenceRepository {
	return &sentenceRepository{db}
}
func (sr *sentenceRepository) GetAllSentences(sentences *[]model.Sentences, userId uint) error {
	if err := sr.db.Joins("User").Where("sentences.user_id=?", userId).Order("created_at").Find(sentences).Error; err != nil {
		return err
	}
	return nil
}

func (sr *sentenceRepository) GetSentenceById(sentence *model.Sentences, userId uint, sentenceId uint) error {
	if err := sr.db.Joins("User").Where("sentences.user_id=?", userId).First(sentence, sentenceId).Error; err != nil {
		return err
	}
	return nil
}

func (sr *sentenceRepository) CreateSentence(sentence *model.Sentences) error {
	if err := sr.db.Create(sentence).Error; err != nil {
		return err
	}
	return nil
}

func (sr *sentenceRepository) DeleteSentence(userId uint, sentenceId uint) error {
	result := sr.db.Where("sentence_id=? AND user_id=?", sentenceId, userId).Delete(&model.Sentences{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}
