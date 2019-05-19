package selection

import (
	"database/sql"
	"regexp"
	"strconv"

	"github.com/rs/zerolog"
)

type DefaultService struct {
	logger     zerolog.Logger
	repository Repository
	sorter     Sorter
	batcher    Batcher
	regex      *regexp.Regexp
}

func NewDefaultService(logger zerolog.Logger, repository Repository, sorter Sorter, batcher Batcher) Service {
	regex := regexp.MustCompile("[0-9]+")
	return &DefaultService{logger, repository, sorter, batcher, regex}
}

func (s DefaultService) Create(req CreateSelectionRequest) (SelectionReply, error) {
	selection, err := s.repository.Selection(req.AppId, req.InstanceId, req.UserId, req.ServerId)
	if err == nil {
		s.logger.Info().
			EmbedObject(selection).
			Msg("found existing selection")

		return s.createSelectionReply(req, selection), nil
	}
	if err != nil && err != sql.ErrNoRows {
		return SelectionReply{}, err
	}

	selection = Selection{
		AppId:      req.AppId,
		InstanceId: req.InstanceId,
		UserId:     req.UserId,
		ServerId:   req.ServerId,
		Options:    map[int]Option{},
	}

	for i, option := range req.Options {
		selection.Options[i+1] = option
	}

	err = s.repository.CreateSelection(selection)
	if err != nil {
		return SelectionReply{}, err
	}

	s.logger.Info().
		EmbedObject(selection).
		Msg("created selection")

	return s.createSelectionReply(req, selection), nil
}

func (s DefaultService) Parse(req ParseSelectionRequest) ([]RankedOption, error) {
	selection, err := s.repository.Selection(req.AppId, req.InstanceId, req.UserId, req.ServerId)
	if err != nil {
		return nil, err
	}

	choices := s.regex.FindAllString(req.Content, -1)

	rankedOptions := []RankedOption{}

	for i, choice := range choices {
		c, err := strconv.Atoi(choice)
		if err != nil {
			return nil, NewValidationError("%s is not a valid integer. %s", choice, err)
		}

		option, ok := selection.Options[c]
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

func (s DefaultService) createSelectionReply(req CreateSelectionRequest, selection Selection) SelectionReply {
	batchOptions := s.createBatchOptions(req, selection)

	sortedBatchOptions := s.sorter.Sort(batchOptions, req.SortMethod, req.SortKey)

	selectionReply := SelectionReply{
		Selection: selection,
		Batches:   s.batcher.CreateBatches(sortedBatchOptions, req.BatchSize),
	}

	return selectionReply
}

func (s DefaultService) createBatchOptions(req CreateSelectionRequest, selection Selection) []BatchOption {
	batchOptions := BatchOptions{}

	for k, option := range selection.Options {
		batchOption := BatchOption{
			Number: k,
			Option: option,
		}

		batchOptions = append(batchOptions, batchOption)
	}

	return batchOptions
}
