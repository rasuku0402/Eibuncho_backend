package usecase

import (
	"testToDoRestAPI/model"
	"testToDoRestAPI/repository"
	"testToDoRestAPI/validator"
)

type ISentenceUsecase interface {
	GetAllSentences(userId uint) ([]model.SentenceResponse, error)
	GetSentenceById(userId uint, sentenceId uint) (model.SentenceResponse, error)
	CreateSentence(sentence model.Sentences) (model.SentenceResponse, error)
	DeleteSentence(userId uint, sentenceId uint) error
}

type sentenceUsecase struct {
	sr repository.ISentenceRepository
	sv validator.ISentenceValidator
}

func NewSentenceUsecase(sr repository.ISentenceRepository, sv validator.ISentenceValidator) ISentenceUsecase {
	return &sentenceUsecase{sr, sv}
}

func (su *sentenceUsecase) GetAllSentences(userId uint) ([]model.SentenceResponse, error) {
	sentences := []model.Sentences{}
	if err := su.sr.GetAllSentences(&sentences, userId); err != nil {
		return nil, err
	}
	resSentences := []model.SentenceResponse{}
	for _, v := range sentences {
		s := model.SentenceResponse{
			SentenceID:   v.SentenceID,
			Japanese_sen: v.Japanese_sen,
			English_sen:  v.English_sen,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		}
		resSentences = append(resSentences, s)
	}
	return resSentences, nil
}

func (su *sentenceUsecase) GetSentenceById(userId uint, sentenceId uint) (model.SentenceResponse, error) {
	sentence := model.Sentences{}
	if err := su.sr.GetSentenceById(&sentence, userId, sentenceId); err != nil {
		return model.SentenceResponse{}, err
	}
	resSentence := model.SentenceResponse{
		SentenceID:   sentence.SentenceID,
		Japanese_sen: sentence.Japanese_sen,
		English_sen:  sentence.English_sen,
		CreatedAt:    sentence.CreatedAt,
		UpdatedAt:    sentence.UpdatedAt,
	}
	return resSentence, nil
}

func (su *sentenceUsecase) CreateSentence(sentence model.Sentences) (model.SentenceResponse, error) {
	if err := su.sv.SentenceValidate(sentence); err != nil {
		return model.SentenceResponse{}, err
	}

	if err := su.sr.CreateSentence(&sentence); err != nil {
		return model.SentenceResponse{}, err
	}
	resSentence := model.SentenceResponse{
		SentenceID:   sentence.SentenceID,
		Japanese_sen: sentence.Japanese_sen,
		English_sen:  sentence.English_sen,
		CreatedAt:    sentence.CreatedAt,
		UpdatedAt:    sentence.UpdatedAt,
	}
	return resSentence, nil
}

func (su *sentenceUsecase) DeleteSentence(userId uint, sentenceId uint) error {
	if err := su.sr.DeleteSentence(userId, sentenceId); err != nil {
		return err
	}
	return nil
}
