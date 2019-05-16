package selection

import (
	"database/sql"
	"math/rand"
	"regexp"
	"sort"
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
	selection, err := s.repository.Selection(req.AppId, req.InstanceId, req.UserId, req.ServerId)
	if err == nil {
		return selection, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return Selection{}, err
	}

	selection = Selection{
		AppId:      req.AppId,
		InstanceId: req.InstanceId,
		UserId:     req.UserId,
		ServerId:   req.ServerId,
		Batches:    []Batch{},
	}

	sortedOptions := s.sortOptions(req.Options, req.SortMethod)

	selection.Batches = s.buildBatches(sortedOptions, req.BatchSize)

	err = s.repository.CreateSelection(selection)
	if err != nil {
		return Selection{}, err
	}

	return selection, nil
}

func (s DefaultService) Parse(req ParseSelectionRequest) ([]RankedOption, error) {
	selection, err := s.repository.Selection(req.AppId, req.InstanceId, req.UserId, req.ServerId)
	if err != nil {
		return nil, err
	}

	options := map[int]Option{}
	for _, batch := range selection.Batches {
		for k, option := range batch.Options {
			options[k] = option
		}
	}

	choices := s.regex.FindAllString(req.Content, -1)

	rankedOptions := []RankedOption{}

	for i, choice := range choices {
		c, err := strconv.Atoi(choice)
		if err != nil {
			return nil, NewValidationError("%s is not a valid integer. %s", choice, err)
		}

		option, ok := options[c]
		if !ok {
			return nil, NewValidationError("could not find option for id: %d", c)
		}

		rankedOption := RankedOption{
			Rank:   i,
			Option: option,
		}

		rankedOptions = append(rankedOptions, rankedOption)
	}

	return rankedOptions, nil
}

func (s DefaultService) sortOptions(options []Option, method SortMethod) []Option {
	switch method {
	case Random:
		return s.shuffleOptions(options)
	case Alphabetical:
		sort.Sort(ByAlphabetical(options))
		return options
	}

	return options
}

func (s DefaultService) shuffleOptions(options []Option) []Option {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := range options {
		j := r.Intn(i + 1)
		options[i], options[j] = options[j], options[i]
	}

	return options
}

func (s DefaultService) buildBatches(options []Option, batchSize int) []Batch {
	numOptions := len(options)

	if batchSize == 0 {
		batchSize = numOptions
	}

	batches := []Batch{}

	for i := 0; i < numOptions; i += batchSize {
		nextBound := i + batchSize

		if nextBound > numOptions {
			nextBound = numOptions
		}

		batchOptions := options[i:nextBound]

		batch := Batch{
			Start:   i + 1,
			End:     nextBound,
			Options: map[int]Option{},
		}

		for j, option := range batchOptions {
			batch.Options[i+j+1] = option
		}

		batches = append(batches, batch)
	}

	return batches
}
