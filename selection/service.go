package selection

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
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
	regex      *regexp.Regexp
}

func NewDefaultService(logger zerolog.Logger, repository Repository) Service {
	regex := regexp.MustCompile("[0-9]+")
	return &DefaultService{logger, repository, regex}
}

func (s DefaultService) Create(req CreateSelectionRequest) (Selection, error) {
	selection := Selection{
		AppId:    req.AppId,
		UserId:   req.UserId,
		ServerId: req.ServerId,
		Options:  map[int]Option{},
	}

	if req.Randomize {
		req.Options = shuffleOptions(req.Options)
	}

	for i, option := range req.Options {
		selection.Options[i+1] = option
	}

	err := s.repository.CreateSelection(selection)
	if err != nil {
		return Selection{}, err
	}

	return selection, nil
}

func (s DefaultService) Parse(req ParseSelectionRequest) ([]RankedOption, error) {
	selection, err := s.repository.Selection(req.AppId, req.UserId, req.ServerId)
	if err != nil {
		return nil, err
	}

	choices := s.regex.FindAllString(req.Content, -1)

	rankedOptions := []RankedOption{}

	for i, choice := range choices {
		c, err := strconv.Atoi(choice)
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid integer. %s", choice, err)
		}

		option, ok := selection.Options[c]
		if !ok {
			return nil, fmt.Errorf("could not find option for id: %d", c)
		}

		rankedOption := RankedOption{
			Rank:   i,
			Option: option,
		}

		rankedOptions = append(rankedOptions, rankedOption)
	}

	return rankedOptions, nil
}

func shuffleOptions(options []Option) []Option {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := range options {
		j := r.Intn(i + 1)
		options[i], options[j] = options[j], options[i]
	}

	return options
}
