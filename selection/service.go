package selection

import (
	"math/rand"
	"time"

	"github.com/rs/zerolog"
)

type Service interface {
	Create(CreateSelectionRequest) (Selection, error)
	Parse(ParseSelectionRequest) ([]RankedOption, error)
}

type DefaultService struct {
	logger     zerolog.Logger
	repository Repository
}

func NewDefaultService(logger zerolog.Logger, repository Repository) Service {
	return &DefaultService{logger, repository}
}

func (s DefaultService) Create(req CreateSelectionRequest) (Selection, error) {
	selection := Selection{
		AppId:    req.AppId,
		UserId:   req.UserId,
		ServerId: req.ServerId,
		Options:  []SelectionOption{},
	}

	if req.Randomize {
		req.Options = shuffleOptions(req.Options)
	}

	for i, option := range req.Options {
		selectionOption := SelectionOption{
			SelectionOptionIndex: i + 1,
			Option:               option,
		}

		selection.Options = append(selection.Options, selectionOption)
	}

	err := s.repository.CreateSelection(selection)
	if err != nil {
		return Selection{}, err
	}

	return selection, nil
}

func (s DefaultService) Parse(req ParseSelectionRequest) ([]RankedOption, error) {
	return []RankedOption{}, nil
}

func shuffleOptions(options []Option) []Option {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := range options {
		j := r.Intn(i + 1)
		options[i], options[j] = options[j], options[i]
	}

	return options
}
