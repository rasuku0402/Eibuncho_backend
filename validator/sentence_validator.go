package validator

import (
	"testToDoRestAPI/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ISentenceValidator interface {
	SentenceValidate(sentence model.Sentences) error
}

type sentenceValidator struct{}

func NewSentenceValidator() ISentenceValidator {
	return &sentenceValidator{}
}

func (sv *sentenceValidator) SentenceValidate(sentence model.Sentences) error {
	return validation.ValidateStruct(&sentence,
		validation.Field(
			&sentence.Japanese_sen,
			validation.Required.Error("Japanese_sentence is required"),
		),
		validation.Field(
			&sentence.English_sen,
			validation.Required.Error("English_sentence is required"),
		),
	)
}
